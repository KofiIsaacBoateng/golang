package main

import (
	"fmt"
	"math"
	"strings"
);


func main() {
	// var declaration
	var name string;
	name = "Kofi Boateng";

	// variable declaration and initialization
	var age int = 24;

	fmt.Println("Name: ", name, " Age: ", age)

	var price float64 = 65.32 // floating 
	var isStocked = false

	fmt.Println("Price: ", price, " Stocked: ", isStocked)


	// short declaration := type inference 
	totalRating, numberOfRating := 1245.5, 1001

	fmt.Println("Average rating: ", totalRating / float64(numberOfRating))
	fmt.Printf("Square root of %T %d is %.1f",8, 8, math.Sqrt(8))


	//String manipulations
	firstName := "james"
	lastName := "bran"
	fullName := firstName + " " + lastName

	fmt.Println("\nMy name (capitalised) is ", strings.ToTitle(fullName))


	//  Integer, Boolean and Floating point operations (usual)

	// constants and typed constants (usual)


	{/**** IF AND ELSE STATEMENTS (also usual minus parenthesis around conditionals) */}

	{/*** For loop (also also usual minus parenthesis) */}

	{/*** Switch statements (also also also usual minus parenthesis and break statements) */}


}