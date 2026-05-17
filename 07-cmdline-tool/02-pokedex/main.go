package main

import (
	"pokedex/internal/pokeapi"
	"time"
)

type Config struct {
	PokeApiClient *pokeapi.Client
	PrevLocationUrl *string
	NextLocationUrl *string
	PokemonCollection map[string]pokeapi.PokemonResp
}

func main() {
	cfg := Config{
		PokeApiClient: pokeapi.NewClient(7 * 24 * time.Hour),
		PokemonCollection: make(map[string]pokeapi.PokemonResp),
	}

	Repl(&cfg);
}