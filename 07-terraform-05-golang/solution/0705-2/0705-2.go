package main

import (
	"fmt"
	"sort"
)

func getMin(a []int) (min int)  {
	min = a[0]
	for i := 0; i < len(a); i++ {
		if min > a[i] {
			min = a[i]
		}
	}

	return
}

func main() {
	x := []int{48, 96, 86, 68, 57, 82, 63, 70, 37, 34, 83, 27, 19, 97, 9, 17}
	fmt.Println("X array:", x)
	fmt.Printf("Minimal value of x array is: %d\n", getMin(x))

	sort.Ints(x)
	fmt.Printf("Truly %d!\n", x[0])
}
