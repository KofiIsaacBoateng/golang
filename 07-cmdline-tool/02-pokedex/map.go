package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func mapCall(cfg *Config) {
	locationAreas, err := cfg.PokeApiClient.LocationAreas(cfg.NextLocationUrl)
	if err != nil {
		log.Fatalf("Location areas: %v", err)
	}

	boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
	boldCyan.Print("\nLocation Areas:\n")

	for _, location := range locationAreas.Results {
		fmt.Printf("[.] %s\n", location.Name)
	}

	fmt.Println("")
	cfg.NextLocationUrl = locationAreas.Next;
	cfg.PrevLocationUrl = locationAreas.Previous;
}

func maprCall(cfg *Config) {
	url := cfg.PrevLocationUrl;
	if url == nil {
		fmt.Println("No need to look back... There is nothing there to find!")
		return
	}
	locationAreas, err := cfg.PokeApiClient.LocationAreas(url);
	if err != nil {
		log.Fatalf("Location areas: %v", err)
	}

	boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
	boldCyan.Print("\nLocation Areas:\n")

	for _, location := range locationAreas.Results {
		fmt.Printf("[.] %s\n", location.Name)
	}

	fmt.Println("")
	cfg.NextLocationUrl = locationAreas.Next;
	cfg.PrevLocationUrl = locationAreas.Previous;
}