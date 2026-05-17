package main

import (
	"fmt"

	"github.com/fatih/color"
)

func mapCall(cfg *Config, args ...string) {
	locationAreas, err := cfg.PokeApiClient.LocationAreas(cfg.NextLocationUrl)
	if err != nil {
		fmt.Printf("Location areas: %v", err)
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

func maprCall(cfg *Config, args ...string) {
	url := cfg.PrevLocationUrl;
	if url == nil {
		fmt.Println("No need to look back... There is nothing there to find!")
		return
	}
	locationAreas, err := cfg.PokeApiClient.LocationAreas(url);
	if err != nil {
		fmt.Printf("Location areas: %v", err)
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