package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	u1 := User{
		Name: "James Bran",
	}

	u1.SayHi()
}

// this method gets a copy of the user
func (u User) SayHi() {
	fmt.Println("Hi there! My name is ", u.Name, ".")
}