package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	u1 := User{
		Name: "James Bran",
		Age: 23,
	}

	fmt.Println("Heyy, ", u1.Name, ". You are ", u1.Age , " years old.")
	u1.Birthday()
	fmt.Println("No. It's my birthday today. I am ", u1.Age, " now!")
}

// this method gives you access the memory val of the User Struct. For modification.
func (u *User) Birthday() {
	u.Age++
}