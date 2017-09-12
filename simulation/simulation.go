package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/drix00/goElectron/simulation/input"
	"github.com/drix00/goElectron/simulation/results"
	"github.com/drix00/goElectron/utils/display"
	"github.com/drix00/goElectron/simulation/models/energyloss"
)

const (
	// Cutoff energy in keV in bulk case.
	cutoffEnergy_keV = 0.5
)

var (
	al_a  float64
	bkSct int
	cp    float64
	er    float64
	ga    float64
	lam_a float64
	sp    float64
	sg_a  float64
	ca    float64
	cb    float64
	cc    float64
	cx    float64
	cy    float64
	cz    float64
	x     float64
	y     float64
	z     float64
	xn    float64
	yn    float64
	zn    float64
	s_en  float64
	dt    display.DisplayTrajectory
)

func ComputeSimulation(inputData input.InputData) results.ResultsData {
	fmt.Println("Single scattering Monte Carlo simulation")
	fmt.Println("Beam energy in keV", inputData.ProbeEnergy_keV)
	fmt.Println("Target atomic number", inputData.AtomicNumber)
	fmt.Println("Target atomic weight", inputData.AtomicWeight_g_mol)
	fmt.Println("Target density_g_cm3 in g/cc", inputData.Density_g_cm3)
	fmt.Println("Foil thickness_A (A)", inputData.Thickness_A)
	fmt.Println("Number of trajectories required", inputData.NumberTrajectories)
	electronRange_A := 700.0 * math.Pow(inputData.ProbeEnergy_keV, 1.66) / inputData.Density_g_cm3
	fmt.Println("Electron range estimate (A)", electronRange_A)

	var results results.ResultsData

	// and parameters needed later.
	getConstants(inputData)

	// Reset the random number generator.
	randomize()

	// Set up the graphic display for plotting.
	initialize()
	setUpGraphics()

	// The Monte Carlo loop.
	zeroCounters()

	imageSize := 700.0 * math.Pow(inputData.ProbeEnergy_keV, 1.66) / inputData.Density_g_cm3 * 1.5
	dt = display.NewDisplayTrajectory(512, imageSize)

	for trajectoryID := 1; trajectoryID < inputData.NumberTrajectories; trajectoryID++ {
		resetCoordinates(inputData)

		// Allow initial entrance of electron.
		step := -computeLambda(s_en) * math.Log(rand.Float64())
		zn = step

		if zn > inputData.Thickness_A {
			// This one is transmitted.
			straightThrough(inputData)
			break
		} else {
			// Plot this position and reset coordinates.
			xyplot(0, 0, 0, zn)
			y = 0.0
			z = zn
		}

		// Start the single scattering loop.

		// Stop when the energy drops below cutoffEnergy_keV.
		for s_en >= cutoffEnergy_keV {
			step = -computeLambda(s_en) * math.Log(rand.Float64())
			sScatter(s_en)
			newCoordinate(step)

			// Problem-specific code will go here.

			// Decide what happens next.
			if zn <= 0.0 {
				// This one is backscattered.
				backscatter()
				break
			}

			if zn > inputData.Thickness_A {
				// This one is transmitted.
				transmit_electron(inputData)
				break
			}

			// Otherwise we go round again.
			resetNextStep(step, inputData)

		}
		// End of the Monte Carlo loop

		if trajectoryID%100 == 0.0 {
			// Reset generator.
			//randomize()
			//showTrajectoryNumbers(trajectoryID, inputData)
		}
	}

	results.BseCoefficient = computeBseCoefficient(inputData)

	filename := fmt.Sprintf("single_scattering_trajectories.png")
	dt.SaveFile(filename)

	return results
}

// Reset the random number generator.
func randomize() {
	seedValue := time.Now().UnixNano()
	rand.Seed(seedValue)
}

// Computes elastic MFP for single scattering model.
func computeLambda(energy float64) float64 {
	al := al_a / energy
	ak := al * (1.0 + al)

	// Giving sg cross section in cm2 as
	sg := sg_a / (energy * energy * ak)

	// and lambda in angstroms is
	return lam_a / sg
}

// Computes some constants needed by the program.
func getConstants(inputData input.InputData) {
	al_a = math.Pow(float64(inputData.AtomicNumber), 0.67) * 3.43e-3

	// Relativistically correct the beam energy for use upt to 500 keV.
	er = (inputData.ProbeEnergy_keV + 511.0) / (inputData.ProbeEnergy_keV + 1022.0)
	er = er * er

	// lambda in cm.
	lam_a = inputData.AtomicWeight_g_mol / (inputData.Density_g_cm3 * 6.0e23)
	// lambda in angstroms.
	lam_a = lam_a * 1.0e8

	sg_a = float64(inputData.AtomicNumber) * float64(inputData.AtomicNumber) * 12.56 * 5.21e-21 * er
}

// Sets up the graphics drivers.
func initialize() {

}

// This displays the trajectories on the pixel screen.
func xyplot(a float64, b float64, c float64, d float64) {
	position := display.Step{a, 0.0, b}
	newPosition := display.Step{c - a, 0.0, d - b}
	dt.DrawStepXZ(position, newPosition)
}

// Draws in the surface(s), beam location, and action thermometers.
func setUpGraphics() {

}

// Resets coordinates at start of each trajectory.
func resetCoordinates(inputData input.InputData) {
	s_en = inputData.ProbeEnergy_keV
	x = 0.0
	y = 0.0
	z = 0.0

	cx = 0.0
	cy = 0.0
	cz = 1.0
}

func zeroCounters() {
	bkSct = 0
}

// Calculates scattering angle using screened Rutherford cross-section.
func sScatter(energy float64) {
	al := al_a / energy
	R1 := rand.Float64()
	cp = 1.0 - ((2.0 * al * R1) / (1.0 + al - R1))
	sp = math.Sqrt(1.0 - cp*cp)

	// and get the azimuthal scattering angle.
	ga = 2.0 * math.Pi * rand.Float64()
}

// Gets xn, yn, zn from x,y,z and scattering angles.
func newCoordinate(step float64) {
	// Find the transformation angles.
	if cz == 0.0 {
		cz = 0.000001
	}
	an_n := (-cx / cz)
	an_m := 1.0 / math.Sqrt(1.0+(an_n*an_n))

	// Save computation time by getting all the transcendentals first.
	v1 := an_m * sp
	v2 := an_m * an_n * sp
	v3 := math.Cos(ga)
	v4 := math.Sin(ga)

	// Find the new direction cosines
	ca = (cx * cp) + (v1 * v3) + (cy * v2 * v4)
	cb = (cy * cp) + (v4 * (cz*v1 - cx*v2))
	cc = (cz * cp) + (v2 * v3) - (cy * v1 * v4)

	//sumC := ca*ca + cb*cb + cc*cc
	//fmt.Println(sumC)

	// and get the new coordinates.
	xn = x + step*ca
	yn = y + step*cb
	zn = z + step*cc
}

// Handles case where initial entry exceeds thickness_A.
func straightThrough(inputData input.InputData) {
	xyplot(0, 0, 0, inputData.Thickness_A)
}

// Handles case of backscattered electron.
func backscatter() {
	bkSct += 1
	xyplot(y, z, yn, zn)
}

// Handles case of transmitted electron.
func transmit_electron(inputData input.InputData) {
	// Length of path from z to bottom face.
	ll := (inputData.Thickness_A - z) / cc
	// Hence the exit y-coordinate.
	yn := y + ll*cb
	xyplot(y, z, yn, inputData.Thickness_A)
}

// Resets variables for next trajectory step.
func resetNextStep(step float64, inputData input.InputData) {
	xyplot(y, z, yn, zn)
	cx = ca
	cy = cb
	cz = cc

	x = xn
	y = yn
	z = zn

	// Find the energy lost on this step.
	delta_E := step * energyloss.StoppingPowerBetheJoyLuo1989(s_en, inputData) * inputData.Density_g_cm3 * 1.0e-8

	// So the current energy is
	s_en = s_en - delta_E
}

// Display number of trajectories done.
func showTrajectoryNumbers(trajectoryID int, inputData input.InputData) {
	fmt.Println("Trajectories done", float64(trajectoryID)/float64(inputData.NumberTrajectories))
}

// Display BSE coefficient on thermometer scale.
func computeBseCoefficient(inputData input.InputData) (float64) {
	bseCoefficient := float64(bkSct) / float64(inputData.NumberTrajectories)
	return bseCoefficient
}
