package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CatFactStruct struct {
	Fact string `json:"fact"`
	Length int `json:"length"`
}

func writeResponseJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func fetchCatFact() (CatFactStruct, error) {
url := "https://catfact.ninja/fact"

	res, err := http.Get(url);

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return CatFactStruct{}, err;
	}

	defer res.Body.Close();

	if res.StatusCode != http.StatusOK {
		fmt.Println("Ops failed with status:", res.Status)
		return CatFactStruct{}, err;
	}

	bodyBytes,err := io.ReadAll(res.Body);

	if err != nil {
		fmt.Println("Error parsing response body to bytes:", err)
		return CatFactStruct{}, err;
	}

	var catfact CatFactStruct;

	if err := json.Unmarshal(bodyBytes, &catfact); err != nil {
		fmt.Println("json unmarshal FAILED!");
		return CatFactStruct{}, err;
	}

	return catfact, nil
}

func externalHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		writeResponseJSON(w, http.StatusInternalServerError, map[string]any {
			"success": false,
			"error": "Invalid Request method. Only GET is allowed!",
		})
		return;
	}

	catfact, err := fetchCatFact();
	if err != nil {
		fmt.Println("Failed to fetch cat data:", err)
		writeResponseJSON(w, http.StatusInternalServerError, map[string]any {
			"success": false,
			"error": err,
			"data": catfact,
		})
		return;
	}

	writeResponseJSON(w, http.StatusOK, map[string]any {
		"success": true,
		"data": map[string]any {
			"source": "Catfact.ninja",
			"fact": catfact.Fact,
			"length": catfact.Length,
		},
		"message": "Cat fact retrieved successfully!",
		"updatedAt": time.Now().UTC(),
	})
}

func main() {
	http.HandleFunc("/external", externalHandler)

	fmt.Println("Server is listening on port: 5000")
	err := http.ListenAndServe(":5000", nil)

	fmt.Println("Listener error:", err)
}