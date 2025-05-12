package math_test

import (
	"basic-unit-testing/math"
	"testing"
)

var mathUtils math.MathUtils = *math.New()

func TestAddition(t *testing.T) {
	res := mathUtils.Add(5, 2)

	if res != 7 {
		t.Errorf("Expected %v, got %v", 7, res)
	}
}
