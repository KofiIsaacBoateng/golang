package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CatFactStruct struct {
	Fact string `json:"fact"`
	Length int `json:"length"`
}

func main() {
	url := "https://catfact.ninja/fact"

	res, err := http.Get(url);

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return;
	}

	defer res.Body.Close();

	if res.StatusCode != http.StatusOK {
		fmt.Println("Ops failed with status:", res.Status)
		return
	}

	bodyBytes,err := io.ReadAll(res.Body);

	if err != nil {
		fmt.Println("Error parsing response body to bytes:", err)
		return
	}

	var catfact CatFactStruct;

	if err := json.Unmarshal(bodyBytes, &catfact); err != nil {
		fmt.Println("json unmarshal FAILED!");
		return;
	}

	fmt.Println("Cat fact:", catfact.Fact)
	fmt.Println("Fact length:", catfact.Length)

}