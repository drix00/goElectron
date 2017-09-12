package energyloss

import (
	"testing"
	"github.com/drix00/goElectron/simulation/input"
)

func BenchmarkStoppingPowerWithBetheJoyLuo1989(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPower(energy_keV, inputData, BetheJoyLuo1989)
	}
}

func BenchmarkStoppingPowerWithBethe1930(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPower(energy_keV, inputData, Bethe1930)
	}
}

func BenchmarkStoppingPowerWithBetheRaoSahibWittry1974(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPower(energy_keV, inputData, BetheRaoSahibWittry1974)
	}
}

func BenchmarkStoppingPowerBetheJoyLuo1989(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPowerBetheJoyLuo1989(energy_keV, inputData)
	}
}

func BenchmarkStoppingPowerBethe1930(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPowerBethe1930(energy_keV, inputData)
	}
}

func BenchmarkStoppingPowerBetheRaoSahibWittry1974(b *testing.B) {
	inputData := input.InputData{6, 12.011, 2.62, 10.0, 1.0e10, 10000, 0.0}
	energy_keV := 20.0
	b.ResetTimer()

	for i :=0 ; i<b.N; i++ {
		StoppingPowerBetheRaoSahibWittry1974(energy_keV, inputData)
	}
}

