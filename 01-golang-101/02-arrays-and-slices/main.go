package main

import (
	"fmt"
)



func main() {
	// arrays have fixed length and cannot grow
	var numbers = [5]int{ 19, 21, 32, 43, 54 }

	numbers[2] = 300

	{/*** 
		* Indexed indexes
		* specific index intitializtion
		* inferred arrays	
	*/}



	{/*** SLICES */}
		// variable sized
		users := []string{"james"}


		var nums []int;

		nums = append(nums, 12)
		nums = append(nums, 23, 34, 54, 54)


		/**** len and capacity */
		// len as usual
		// capacity is how many elements a slice can hold
		scores := make([]int, 0,10) // variable len, capacity of 10
		scores = append(scores, 12, 23, 34, 34, 34, 34, 3,4, 3,34, 34, 34, 34)

		// spreading is as usual but inverse eg. users... instead of ...users
		// FOR RANGE is javascripts forEach. for index, view := range views {}

	fmt.Println(numbers)
	fmt.Println(users)
	fmt.Println(nums, "\nLength of nums: ", len(nums))
	fmt.Println("Length: ", len(scores), ", capacity: ", cap(scores), ", scores: ", scores)
}