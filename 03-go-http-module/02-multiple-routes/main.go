package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Welcome to the main route! Try '/hell0?name=< >' with your name between the angle brackets."))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == ""{
		name = "Guest"
	}

	_,_ = w.Write([]byte(fmt.Sprintf("Hello %s. You made it!", name)))
}


func main(){
	// root route
	http.HandleFunc("/", rootHandler)
	// hello route
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server is listening on port: 5000")
	err := http.ListenAndServe(":5000", nil)

	fmt.Println("Listener error:", err)
}