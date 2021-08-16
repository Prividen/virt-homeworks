package main

import "fmt"

func m2f (meters float64) (foots float64) {
	foots = meters / 0.3048
	return
}

func main() {
	var meters float64
	fmt.Print("Please enter meters value: ")
	_, err := fmt.Scanf("%f", &meters)

	if err != nil {
		fmt.Println("Please specify correct meters value")
		return
	}

	fmt.Printf("%f meters is %f foots\n", meters, m2f(meters))
}
