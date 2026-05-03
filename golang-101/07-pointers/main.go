package main

import "fmt"

func increaseScore(score *int, by int) {
	*score += by
}

func decreaseScore(score *int, by int) {
	*score -= by
}

func main() {
	// inference(&var) => access the address of the variable var
	// deference(*var) => access the value saved at the address &var
	
	score := 20

	increaseScore(&score, 5)
	decreaseScore(&score, 2)
	decreaseScore(&score, 2)
	increaseScore(&score, 5)
	increaseScore(&score, 5)
	decreaseScore(&score, 3)
	increaseScore(&score, 5)

	fmt.Println("Final score: ", score)
}