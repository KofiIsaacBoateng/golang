package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)


func searchFile (file string, query string){
	content, err := os.ReadFile(file);
	if(err != nil) {
		fmt.Printf("Failed to read %s: %v\n", file, err);
		return;
	}

	matched := strings.Contains(string(content), query);
	if matched {
		fmt.Printf("[FOUND] Match found in file %s\n", file);
	}else {
		fmt.Printf("[NOT FOUND] No match found in %s\n", file)
	}
}

func main() {
	query := "concurrency"
	files := []string{
		"./lesson/files/file0.txt",
		"./lesson/files/file1.txt",
		"./lesson/files/file2.txt",
		"./lesson/files/file3.txt",
	}

	fmt.Print("\n[.] Starting concurrency search...\n\n")

	// Launches a separate goroutine for each file with the go keyword
	for _, file := range files {
		
		go searchFile(file, query)

	}

	// Without this, the main goroutine exits instantly before the others finish.
	time.Sleep(100 * time.Millisecond);

	fmt.Println("\n\n[.] Search operation finished.")
}