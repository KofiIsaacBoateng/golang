package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only get requests are allowed!", http.StatusMethodNotAllowed);
		return
	}

	_, _ = w.Write([]byte("Hello there! This is my first api server in GO!"))
}

func main() {
	// route handler
	http.HandleFunc("/hello", helloHandler)

	// port listener
	fmt.Println("Server is listening on PORT: 5000")
	err := http.ListenAndServe(":5000", nil)

fmt.Println("error from listener:", err)
}	