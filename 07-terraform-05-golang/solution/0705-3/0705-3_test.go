package main
import (
	"testing"
)

func TestGetDivisibleDigits (t *testing.T) {
	divisibleDigits := getDivisibleDigits(100, 3)
	if len(divisibleDigits) != 33 {
		t.Error("Expected len array: 33")
	}
	if divisibleDigits[0] != 3 {
		t.Error("Expected first element: 3")
	}
	if divisibleDigits[len(divisibleDigits)-1] != 99 {
		t.Error("Expected last element: 99")
	}
}
