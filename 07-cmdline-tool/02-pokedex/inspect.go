package main

import (
	"fmt"

	"github.com/fatih/color"
)



func inspectCall(cfg *Config, args ...string) {
	if len(args) != 1 {
		fmt.Println("Invalid pokemon name!")
		return
	}

	pokemonName := args[0]

	pokemon, ok := cfg.PokemonCollection[pokemonName];
	if !ok {
		fmt.Println("You haven't caught this pokemon yet!")
	}else {
		boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
		// boldCyan.Printf("\n%s:\n", pokemon.Name)

		fmt.Printf("\n[.] Name: %s\n", pokemon.Name)
		fmt.Printf("[.] Height: %v\n", pokemon.Height)
		fmt.Printf("[.] Weight: %v\n", pokemon.Weight)
		boldCyan.Println("[.] stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("   [.] %s:\t%v\n", stat.Stat.Name, stat.BaseStat)
		}

	}
	fmt.Println("")
}
