// Tests for the main file.
package main

import (
	"testing"
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