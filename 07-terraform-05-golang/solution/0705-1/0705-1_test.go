package main

import (
	"testing"
)

func Test_m2f(t *testing.T) {
	testFoots := m2f(0.3048)

	if testFoots != 1 {
		t.Error("Wrong foot value")
	}
}
