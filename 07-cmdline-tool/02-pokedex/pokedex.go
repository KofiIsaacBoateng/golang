package main

import (
	"fmt"

	"github.com/fatih/color"
)



func pokedexCall(cfg *Config, args ...string) {
	boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
	boldCyan.Println("\nYour Pokemon Collection")

	for key := range cfg.PokemonCollection {
		fmt.Printf("[.] %s\n", key)
	}
	
	fmt.Println("")
}
