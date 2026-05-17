package main

import (
	"fmt"

	"github.com/fatih/color"
)



func exploreCall(cfg *Config, args ...string) {
	if len(args) != 1 {
		fmt.Println("Invalid location area name!")
		return
	}

	location := args[0]
	locationArea, err := cfg.PokeApiClient.ExploreLocationArea(location)
	if err != nil {
		fmt.Printf("Error finding Pokemons in this Area: %s => %v", location, err)
	}

	boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
	boldCyan.Printf("\nPokemons in Area: %s\n", location)

	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf("[.] %s\n", pokemon.Pokemon.Name)
	}

	fmt.Println("")
}
