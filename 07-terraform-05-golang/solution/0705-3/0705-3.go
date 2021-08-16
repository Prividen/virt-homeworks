package main

import "fmt"

func getDivisibleDigits (max int, divider int) ([]int) {
	var divisible_digits []int
	for i:=1; i<=max; i++ {
		if float32(i / divider) == float32(i) / float32(divider) {
			divisible_digits = append(divisible_digits, i)
		}
	}
	return divisible_digits
}

func main() {
	fmt.Println(getDivisibleDigits(100, 3))
}
