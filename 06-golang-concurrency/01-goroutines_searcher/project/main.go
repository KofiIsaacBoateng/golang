package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

func searchFile(path string, query string, wg *sync.WaitGroup) {
	// signal goroutine that this wait group is finished when the function exits.
	defer wg.Done()

	// open file
	file, err := os.Open(path);
	if(err != nil) {
		fmt.Printf("Failed to open %s: %v\n", path, err);
		return;
	}

	// Ensure file is closed when the function finishes;
	defer func(){
		if err := file.Close(); err != nil {
			fmt.Printf("File failed to CLOSE: %v", err)
		}
	}()

	// create a scanner to read the file line by line;
	scanner := bufio.NewScanner(file);

	lineCounter := 1;
	result := []string{};
	for scanner.Scan() {
		content := scanner.Text(); // get the content of a line
		matched := strings.Contains(content, query);
		if matched {
			result = append(result, fmt.Sprintf("line [%d] -> %s", lineCounter, content))
		}
		lineCounter++
	}
	
	if len(result) > 0 {
		fmt.Printf("[FOUND] Match found in file %s\n", strings.ReplaceAll(path, `..\..\`, ""));
		for _ , val := range result {
			fmt.Printf("\t%s\n", val);
		}
	}

}


func main () {
	// Search the entire golang journey directory for matches
	// assuming you are in the 01-goroutines_searcher directory
	dir := "../../";
	query := os.Args[1];
	searchIgnorePaths := []string{"git", "node_modules"}; // skip files in these directories := all your .gitignore vals go here


	// declare a waitgroup.
	// it counts how many concurrent tasks are currently running

	var wg sync.WaitGroup;

	// search through the entire tree of dir and find matches in files
	fmt.Printf("Deep scanning scope '%s' for term: '%s'...\n\n", dir, query);

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error)error{
		if err != nil {
			return err;
		}

		// skip directories;
		if info.IsDir() {
			return nil;
		}

		// skip paths marked as ignore
		if shouldIgnore := slices.ContainsFunc(searchIgnorePaths, func(slice string) bool {
			return strings.Contains(path, slice)
		}); shouldIgnore {
			return nil
		}

		// tell wait group that you are spinning a new goroutine
		wg.Add(1)


		// add a concurrency search to this file
		go searchFile(path, query, &wg);

		
		return nil;
	});

	if err != nil {
		fmt.Printf("Error walking directory path %s\n", err)
		return;
	}

	// Block execution here until the counter inside WaitGroup returns to 0
	wg.Wait();

	fmt.Println("\n\n[.] Search Complete!")
}