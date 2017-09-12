// A single scattering Monte Carlo simulation which uses the
// screened Rutherford cross section.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/drix00/goElectron/simulation"
	"github.com/drix00/goElectron/simulation/input"
	"github.com/drix00/goElectron/simulation/models/energyloss"
)

// This is the start of the main program.
func main() {
	inputData := getInputData()

	results := simulation.ComputeSimulation(inputData)

	fmt.Println("BSE coefficient: ", results.BseCoefficient)
}

// Gets the input data to run this program from the command line.
func getInputData() input.InputData {
	var inputData input.InputData

	inputData.ProbeEnergy_keV, _ = strconv.ParseFloat(os.Args[1], 64)
	inputData.AtomicNumber, _ = strconv.Atoi(os.Args[2])
	inputData.MeanIonizationPotential_keV = energyloss.MeanIonizationPotential_keV(inputData.AtomicNumber)
	inputData.AtomicWeight_g_mol, _ = strconv.ParseFloat(os.Args[3], 64)
	inputData.Density_g_cm3, _ = strconv.ParseFloat(os.Args[4], 64)
	inputData.Thickness_A, _ = strconv.ParseFloat(os.Args[5], 64)
	inputData.NumberTrajectories, _ = strconv.Atoi(os.Args[6])

	return inputData
}
