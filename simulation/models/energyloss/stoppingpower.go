package energyloss

import (
	"math"

	"github.com/drix00/goElectron/simulation/input"
)

type StoppingPowerModel uint

const (
	BetheJoyLuo1989         StoppingPowerModel = iota
	Bethe1930               StoppingPowerModel = iota
	BetheRaoSahibWittry1974 StoppingPowerModel = iota
)

// This compute the stopping power in keV/g/cm2.
func StoppingPower(energy float64, inputData input.InputData, model StoppingPowerModel) float64 {
	switch model {
	case BetheJoyLuo1989:
		return StoppingPowerBetheJoyLuo1989(energy, inputData)
	case Bethe1930:
		return StoppingPowerBethe1930(energy, inputData)
	case BetheRaoSahibWittry1974:
		return StoppingPowerBethe1930(energy, inputData)
	default:
		return StoppingPowerBetheJoyLuo1989(energy, inputData)
	}
}

// This compute the stopping power in keV/g/cm2 using the
// modified Bethe expression of Eq. (3.21).
func StoppingPowerBetheJoyLuo1989(energy float64, inputData input.InputData) float64 {
	if energy < 0.05 {
		// Just in case.
		energy = 0.05
	}

	temp := math.Log(1.166 * (energy + 0.85*inputData.MeanIonizationPotential_keV) / inputData.MeanIonizationPotential_keV)
	return temp * 78500.0 * float64(inputData.AtomicNumber) / (inputData.AtomicWeight_g_mol * energy)
}

// This compute the stopping power in keV/g/cm2 using the
// Bethe expression of Eq. (3.17).
func StoppingPowerBethe1930(energy float64, inputData input.InputData) float64 {
	if energy < 0.05 {
		// Just in case.
		energy = 0.05
	}

	temp := math.Log(1.166 * energy / inputData.MeanIonizationPotential_keV)
	return temp * 78500.0 * float64(inputData.AtomicNumber) / (inputData.AtomicWeight_g_mol * energy)
}

// This compute the stopping power in keV/g/cm2 using the
// Bethe expression of Eq. (3.17).
func StoppingPowerBetheRaoSahibWittry1974(energy float64, inputData input.InputData) float64 {
	if energy < 6.4*inputData.MeanIonizationPotential_keV {
		nominator := 62400.0 * float64(inputData.AtomicNumber)
		denominator := math.Sqrt(energy * inputData.MeanIonizationPotential_keV * float64(inputData.AtomicNumber))
		return nominator / denominator
	} else {
		return StoppingPowerBethe1930(energy, inputData)
	}
}
