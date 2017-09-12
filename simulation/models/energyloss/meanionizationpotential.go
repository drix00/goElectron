package energyloss

import "math"

// Calculate J the mean ionization potential meanIonizationPotential_keV using the
// Berger-Seltzer analytical fit.
func MeanIonizationPotential_keV(atomicNumber int) (float64) {
	potential_keV := (9.76*float64(atomicNumber) + (58.5 / math.Pow(float64(atomicNumber), 0.19))) * 0.001
	return potential_keV
}
