package simulation

import (
	"fmt"
	"testing"

	"github.com/drix00/goElectron/simulation/input"
	"github.com/drix00/goElectron/simulation/models/energyloss"
)

func TestZeroCounters(t *testing.T) {
	t.Log("Test if the zeroCounters function reset the counters.")
	bkSct = 3

	if bkSct == 0 {
		t.Error("The variable bkSCt was not set correctly.")
	}

	zeroCounters()

	if bkSct == 0 {
		t.Log("The variable bkSCt was reset correctly.")
	} else {
		t.Error("The variable bkSCt was not reset correctly.")
	}
}

func TestCarbon10keV(t *testing.T) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	inputData.MeanIonizationPotential_keV = energyloss.MeanIonizationPotential_keV(inputData.AtomicNumber)

	results := ComputeSimulation(inputData)

	valueString := fmt.Sprintf("%.2f", results.BseCoefficient)
	if valueString != "0.06" {
		t.Fatalf("BSE coefficient is not correct (0.0634), got %f", results.BseCoefficient)
	}
}
