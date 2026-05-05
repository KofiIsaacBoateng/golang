package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func successHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]any{
		"ok": true,
		"message": "JSON encoded successfully.",
		"updatedAt": time.Now().UTC(), 
	}

	_ = json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/ok", successHandler)

	fmt.Println("Server is listening on port: 5000")
	err := http.ListenAndServe(":5000", nil)

	fmt.Println("Listener error:", err)
}