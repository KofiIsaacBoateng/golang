package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/fatih/color"
)



func catchCall(cfg *Config, args ...string) {
	if len(args) != 1 {
		fmt.Println("Invalid pokemon name!")
		return
	}

	pokemonName := args[0]	
	pokemon, err := cfg.PokeApiClient.CatchPokemon(pokemonName)
	if err != nil {
		fmt.Printf("Error finding Pokemon: %s => %v", pokemonName, err)
	}

	boldCyan := color.New(color.FgCyan, color.Bold, color.Underline);
	boldCyan.Printf("\n%s: ", pokemonName)

	const threshold = 40;
	if(rand.Int64N(pokemon.BaseExperience) > threshold){
		fmt.Println("Yeaaaaaah! I'm untouchable Mother******!!!")
	}else {
		fmt.Println("Congratulations! You caught me!")
		cfg.PokemonCollection[pokemonName] = pokemon;
	}

	fmt.Println("")
}
