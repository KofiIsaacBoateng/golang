package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Age   int
	Email string
}

func main() {
	u1 := User{
		Name:  "Kofi Boateng",
		Age:   32,
		Email: "k@gmail.com",
		ID:    1,
	}

	// mutable by default
	u1.Age = 24

	fmt.Printf("Hi %s. You are %d years old. Your email is %s, and you are #%d.\n", u1.Name, u1.Age, u1.Email, u1.ID)


	// partial User
	u2 := User{
		Name: "William Bran",
		Age: 54,
	}
	fmt.Printf("Hi %s. You are %d years old.\n", u2.Name, u2.Age)

}