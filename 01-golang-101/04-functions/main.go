package main

import (
	"fmt"
)



func add(a int, b int) int {
	return a + b;
}


func prod (a int, b int) int {
	return a * b;
}


func addAndProd (a int, b int) (int, int) {
	return prod(a, b), add(a, b)
}


/**** Named return values */
func subAndDivide(a int, b int) (sub int, divide int) {
	sub = a - b;
	divide = a / b;

	return
}


func main() {
	a, b := 10, 15;
	fmt.Printf("The product and sum of %d and %d are %s\n", a, b, fmt.Sprint(addAndProd(a, b)))
	fmt.Printf("The difference and integer division of  %d and %d are %v.\n", a, b, fmt.Sprint(subAndDivide(a, b)))

	// IIFE
	logger := func(a int, b int) string {
		return fmt.Sprintf("The numbers received are %d and %d", a, b)
	}(12, 15)

	fmt.Println(logger)
}