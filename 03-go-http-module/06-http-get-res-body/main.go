package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://jsonplaceholder.typicode.com/todos"

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

	bodyText := string(bodyBytes);

	max := 250;

	if len(bodyText) > max {
		max = len(bodyText)
	}

	fmt.Println(bodyText)
	fmt.Println("max-length:", max)
}