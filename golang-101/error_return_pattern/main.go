package main

import (
	"fmt"
	"log"
	"strconv"
)


func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}
}


func run() error {
	input := "2"

	level, err := parseLevel(input)

	if err != nil {
		return err
	}

	fmt.Println("You've selected level ", level)
	return nil
}


func parseLevel(s string)(int, error) {

	n, err := strconv.Atoi(s)

	if err != nil {
		return 0, fmt.Errorf("Level must be a number!")
	}

	if(n < 1 || n > 5){
		return 0, fmt.Errorf("Invalid level number. Accepted ranges are 1 to 5.")
	}

	return n, nil
}