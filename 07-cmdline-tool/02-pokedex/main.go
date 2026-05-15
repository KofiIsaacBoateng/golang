package main

import (
	"pokedex/internal/pokeapi"
	"time"
)

type Config struct {
	PokeApiClient *pokeapi.Client
	PrevLocationUrl *string
	NextLocationUrl *string
}

func main() {
	cfg := Config{
		PokeApiClient: pokeapi.NewClient(7 * 24 * time.Hour),
	}

	Repl(&cfg);
}