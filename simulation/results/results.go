package results

import (
	"math"
)

type ResultsData struct {
	BseCoefficient float64
}

// Compute the statistical error on the BSE coefficient from the number of trajectories.
func (rd ResultsData) GetBseError(numberTrajectories int) float64 {
	error := math.Sqrt(rd.BseCoefficient * (1.0 - rd.BseCoefficient) / float64(numberTrajectories))
	return error
}
