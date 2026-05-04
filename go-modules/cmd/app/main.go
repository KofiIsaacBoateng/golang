package main

import (
	"fmt"
	"go-modules/internal/greet"
)

func main() {

	fmt.Println("Case 1: Real user")
	fmt.Println(greet.Hello(" Kofi boateng "))

	fmt.Println("Case 2: Guest user ")
	fmt.Println(greet.Hello("  "))
}