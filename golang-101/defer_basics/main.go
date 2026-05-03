package main

import (
	"errors"
	"fmt"
)

func main() {
	// do Something

	// case 1: Success
	fmt.Println("Case 1: Success.")
	if err := doSomething(true); err != nil {
		fmt.Println("error: ", err)
	}

	// case 2: Failure
	fmt.Println("Case 2: Failure.")
	if err := doSomething(false); err != nil {
		fmt.Println("error: ", err)
	}
}

func doSomething(success bool) error {
	// acquire resourse
	// do work
	// clean up resourses everytime.

	fmt.Println("Resource has been acquired!")

	defer fmt.Println("Cleanup: Clearing out resources.")

	if(!success) {
		return errors.New("Something went wrong. Shutting system down early...")
	}

	// do work
	fmt.Println("Doing some really important work with the resources obtained!")
	fmt.Println("Work completed successfully!")
	return nil
}