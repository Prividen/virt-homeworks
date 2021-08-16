package main

import (
	"testing"
)

func TestGetMin(t *testing.T) {
	min := getMin([]int{8, 2, 15})
	if min != 2 {
		t.Error("Wrong min value")
	}
}

