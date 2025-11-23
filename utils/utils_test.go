package utils

import "testing"

// TestNoUtilFunction tests the correctness of NoUtilFunction by checking if it returns the expected value and error.
func TestNoUtilFunction(t *testing.T) {
	value, err := NoUtilFunction()
	if value != 0 || err != nil {
		t.Errorf("Expected: 0, <nil> || Got: %v, %v", value, err)
	}
}
