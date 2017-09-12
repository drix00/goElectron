// Test the mean ionization potential model.
package energyloss

import "testing"

// Test with model and date from Joy Monte Carlo book.
func TestMeanIonizationPotential_keV(t *testing.T) {
	var values = []struct {
		Z int
		J_keV float64
	} {
		{1, 0.068260,},
		{3, 0.076759,},
		{14, 0.172072,},
		{79, 0.796544,},
	}

	for _, valueRef := range values {
		computedJ_keV := MeanIonizationPotential_keV(valueRef.Z)

		if !floatEquals(valueRef.J_keV, computedJ_keV, 1.0e-5) {
			t.Errorf("Mean ionization potential calculation is wrong: expected %f got %f", valueRef.J_keV, computedJ_keV)
		}
	}
}

func floatEquals(a float64, b float64, epsilon float64) bool {
	if (a - b) < epsilon && (b - a) < epsilon {
		return true
	}
	return false
}