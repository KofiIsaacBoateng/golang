package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func writeResponseJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}

type ReqBodyStruct struct {
	Name string `json:"name"`
}

func jsonDecoderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error": "Invalid request method! Only POST is allowed!",
		})
		return;
	}

	defer r.Body.Close()

	var req ReqBodyStruct;

	dec := json.NewDecoder(r.Body);

	if err := dec.Decode(&req); err != nil {
		writeResponseJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error": "Invalid JSON format!",
		})

		return;
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeResponseJSON(w, http.StatusBadRequest, map[string]any {
			"success": false,
			"error": "Name cannot be an empty string!",
		})
		return;
	}

	writeResponseJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": req,
		"message": "JSON decoded successfully!",
		"updatedAt": time.Now().UTC(),
	})
}


func main(){
	http.HandleFunc("/json-decoder", jsonDecoderHandler)

	fmt.Println("Server is listening on port: 5000")
	err := http.ListenAndServe(":5000", nil)

	fmt.Println("Listener error:", err)
}