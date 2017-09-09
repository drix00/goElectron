// A single scattering Monte Carlo simulation which uses the
// screened Rutherford cross section.
package main

import (
	"math"
	"math/rand"
	"time"
	"fmt"
	"os"
	"strconv"

	"github.com/goElectron/utils/display"
)

const (
	// Cutoff energy in keV in bulk case.
	e_min = 0.5
)

var (
	at_num int
	at_wht float64
	density float64
	inc_energy float64
	mn_ion_pot float64
	al_a float64
	bkSct int
	cp float64
	er float64
	ga float64
	lam_a float64
	sp float64
	sg_a float64
	ca float64
	cb float64
	cc float64
	cx float64
	cy float64
	cz float64
	x float64
	y float64
	z float64
	xn float64
	yn float64
	zn float64
	s_en float64
	thickness float64
	numberTrajectories int
	dt display.DisplayTrajectory
)

// This is the start of the main program.
func main() {
	// Get input data and find J value
	setUpScreen()

	// and parameters needed later.
	getConstants()

	// Reset the random number generator.
	randomize()

	// Set up the graphic display for plotting.
	initialize()
	setUpGraphics()

	// The Monte Carlo loop.
	zeroCounters()

	imageSize := 700.0 * math.Pow(inc_energy, 1.66) / density * 1.5
	dt = display.NewDisplayTrajectory(512, imageSize)

	for trajectoryID :=1; trajectoryID<numberTrajectories; trajectoryID++ {
		resetCoordinates()


		// Allow initial entrance of electron.
		step := -computeLambda(s_en) * math.Log(rand.Float64())
		zn = step


		if zn > thickness {
			// This one is transmitted.
			straightThrough()
			break
		} else {
			// Plot this position and reset coordinates.
			xyplot(0, 0, 0, zn)
			y = 0.0
			z = zn
		}

		// Start the single scattering loop.

		// Stop when the energy drops below e_min.
		for s_en >= e_min {
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

			if zn > thickness {
				// This one is transmitted.
				transmit_electron()
				break
			}

			// Otherwise we go round again.
			resetNextStep(step)

		}
		// End of the Monte Carlo loop

		if trajectoryID % 100 == 0.0 {
			// Reset generator.
			//randomize()
			showTrajectoryNumbers(trajectoryID)
		}
	}

	showBseCoefficient()

	filename := fmt.Sprintf("single_scattering_trajectories.png")
	dt.SaveFile(filename)
}

// Reset the random number generator.
func randomize() {
	seedValue := time.Now().UnixNano()
	rand.Seed(seedValue)
}

// This compute the stopping power in keV/g/cm2 using the
// modified Bethe expression of Eq. (3.21).
func stop_pwr(energy float64) (float64) {
	if energy < 0.05 {
		// Just in case.
		energy = 0.05
	}

	temp := math.Log(1.166 * (energy + 0.85*mn_ion_pot)/mn_ion_pot)
	return temp * 78500.0* float64(at_num) / (at_wht * energy)
}

// Computes elastic MFP for single scattering model.
func computeLambda(energy float64) (float64) {
	al := al_a / energy
	ak := al * (1.0 + al)

	// Giving sg cross section in cm2 as
	sg := sg_a / (energy*energy*ak)

	// and lambda in angstroms is
	return lam_a/sg
}

// Computes some constants needed by the program.
func getConstants()  {
	al_a = math.Pow(float64(at_num), 0.67) * 3.43e-3

	// Relativistically correct the beam energy for use upt to 500 keV.
	er = (inc_energy + 511.0) / (inc_energy + 1022.0)
	er = er *er

	// lambda in cm.
	lam_a = at_wht / (density*6.0e23)
	// lambda in angstroms.
	lam_a = lam_a *1.0e8

	sg_a = float64(at_num) * float64(at_num) * 12.56 * 5.21e-21 * er
}

// Gets the input data to run this program.
func setUpScreen() {
	fmt.Println("Single scattering Monte Carlo simulation")

	inc_energy, _ = strconv.ParseFloat(os.Args[1], 64)
	fmt.Println("Beam energy in keV", inc_energy)

	at_num, _ = strconv.Atoi(os.Args[2])
	fmt.Println("Target atomic number", at_num)

	// Calculate J the mean ionization potential mn_ion_pot using the
	// Berger-Selzer analytical fit.
	mn_ion_pot = (9.76*float64(at_num) + (58.5 / math.Pow(float64(at_num), 0.19))) * 0.001

	at_wht, _ = strconv.ParseFloat(os.Args[3], 64)
	fmt.Println("Target atomic weight", at_wht)

	density, _ = strconv.ParseFloat(os.Args[4], 64)
	fmt.Println("Target density in g/cc", density)

	if len(os.Args) == 7 {
		thickness, _ = strconv.ParseFloat(os.Args[5], 64)
		fmt.Println("Target is thin film")
		fmt.Println("Foil thickness (A)", thickness)

		numberTrajectories, _ = strconv.Atoi(os.Args[6])
		fmt.Println("Number of trajectories required", numberTrajectories)
	} else {
		// Estimate the beam range.
		thickness = 700.0 * math.Pow(inc_energy, 1.66) / density
		fmt.Println("Target is bulk")
		fmt.Println("Electron range estimate (A)", thickness)

		numberTrajectories, _ = strconv.Atoi(os.Args[5])
		fmt.Println("Number of trajectories required", numberTrajectories)
	}
}

// Sets up the graphics drivers.
func initialize()  {

}

// This displays the trajectories on the pixel screen.
func xyplot(a float64, b float64, c float64, d float64)  {
	position := display.Step{a, 0.0, b}
	newPosition := display.Step{c-a, 0.0, d-b}
	dt.DrawStepXZ(position, newPosition)
}

// Draws in the surface(s), beam location, and action thermometers.
func setUpGraphics()  {

}

// Resets coordinates at start of each trajectory.
func resetCoordinates()  {
	s_en = inc_energy
	x = 0.0
	y = 0.0
	z = 0.0

	cx = 0.0
	cy = 0.0
	cz = 1.0
}

func zeroCounters()  {
	bkSct = 0
}

// Calculates scattering angle using screened Rutherford cross-section.
func sScatter(energy float64)  {
	al := al_a / energy
	R1 := rand.Float64()
	cp = 1.0 - ((2.0*al*R1)/(1.0 + al - R1))
	sp = math.Sqrt(1.0 - cp*cp)

	// and get the azimuthal scattering angle.
	ga = 2.0 *math.Pi * rand.Float64()
}

// Gets xn, yn, zn from x,y,z and scattering angles.
func newCoordinate(step float64)  {
	// Find the transformation angles.
	if cz == 0.0 {
		cz = 0.000001
	}
	an_n := (-cx/cz)
	an_m := 1.0 / math.Sqrt(1.0 + (an_n*an_n))

	// Save computation time by getting all the transcendentals first.
	v1 := an_m * sp
	v2 := an_m * an_n * sp
	v3 := math.Cos(ga)
	v4 := math.Sin(ga)

	// Find the new direction cosines
	ca = (cx * cp) + (v1 * v3) + (cy * v2 * v4)
	cb = (cy * cp) + (v4*(cz*v1 - cx*v2))
	cc = (cz * cp) + (v2*v3) - (cy*v1*v4)

	//sumC := ca*ca + cb*cb + cc*cc
	//fmt.Println(sumC)

	// and get the new coordinates.
	xn = x + step * ca
	yn = y + step * cb
	zn = z + step * cc
}

// Handles case where initial entry exceeds thickness.
func straightThrough()  {
	xyplot(0, 0, 0, thickness)
}

// Handles case of backscattered electron.
func backscatter()  {
	bkSct += 1
	xyplot(y, z, yn, zn)
}

// Handles case of transmitted electron.
func transmit_electron()  {
	// Length of path from z to bottom face.
	ll := (thickness - z) / cc
	// Hence the exit y-coordinate.
	yn := y + ll * cb
	xyplot(y, z, yn, thickness)
}

// Resets variables for next trajectory step.
func resetNextStep(step float64)  {
	xyplot(y, z, yn, zn)
	cx = ca
	cy = cb
	cz = cc

	x = xn
	y = yn
	z = zn

	// Find the energy lost on this step.
	delta_E := step * stop_pwr(s_en) * density *1.0e-8

	// So the current energy is
	s_en = s_en - delta_E
}

// Display number of trajectories done.
func showTrajectoryNumbers(trajectoryID int)  {
	fmt.Println("Trajectories done", float64(trajectoryID)/float64(numberTrajectories))
}

// Display BSE coefficient on thermometer scale.
func showBseCoefficient()  {
	bseCoefficient := float64(bkSct)/float64(numberTrajectories)
	fmt.Println("BSE coefficient: ", bseCoefficient)
}
