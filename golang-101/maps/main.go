package main

import (
	"fmt"
)


func main() {
	// map[keyType]valueType
	var ages map[string]int; // nill map
	ages = map[string]int{
		"james": 18,
		"william": 35,
	}

	values := make(map[string]int) // using make

	values["a"] = 1
	values["b"] = 2
	values["c"] = 3
	values["d"] = 4


	// map operations
	users := map[string]string{
		"user1": "Will",
		"user2": "Martin",
		"user3": "Samantha",
		"user4": "George",
		"user5": "Sandra",
	}

	//delete
	delete(users, "user4") // deleted user 4
	delete(users, "user100") // no user 100. However, throws no error


	// value, ok (comma-ok idiom)
	_, ok := users["user100"]; // more like a redundant version of null or falsy values

	if ok {
		fmt.Println("User 100 exists.")
		
	}else {
		fmt.Println("User 100 doesn't exist!")
	}


	// for range also works for Maps allowing access to key and values, just like Object.entries(obj) in javascript. for key, value := range map{}


	fmt.Println(ages, ages["william"], len(ages))
	fmt.Println("Values: ", values)
	fmt.Println("Users: ", users)

}