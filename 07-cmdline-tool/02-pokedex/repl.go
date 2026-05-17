package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Command struct {
	Name string
	Command string
	Description string
	Call func(cfg *Config, args ...string)
}


func Repl(cfg *Config) {
	for {
		scanner := bufio.NewScanner(os.Stdin);

		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgHiGreen).SprintFunc()

		fmt.Printf("%s %s ", cyan("pokedex"), green("$>"))

		scanner.Scan()
		input := strings.TrimSpace(scanner.Text());
		words := strings.Fields(input)

		// onPressEnter
		if len(words) == 0 {
			continue;
		}


		commands := getCommands();

		command, ok := commands[words[0]];
		if !ok {
			// invalid input
			fmt.Printf("command: %s is INVALID.\n", input)
			continue;
		}
		command.Call(cfg, words[1:]...);
	}
}

func getCommands() map[string]Command {
	return map[string]Command{
		"exit": {
			Name: "Exit",
			Command: "exit",
			Description: "Closes pokedex.",
			Call: exitCall,
		},
		"help": {
			Name: "Help",
			Command: "help",
			Description: "Displays details about how to use pokedex",
			Call: helpCall,
		},
		"map": {
			Name: "Map",
			Command: "map",
			Description: "List the locations where pokemons can be found and caught",
			Call: mapCall,
		},

		"mapr": {
			Name: "Reverse Map",
			Command: "mapr",
			Description: "Reverse to previous pages listed by map.",
			Call: maprCall,
		},

		"explore": {
			Name: "Explore",
			Command: "explore [location]",
			Description: "Explore pokemons available at a location.",
			Call: exploreCall,
		},

		"catch": {
			Name: "Catch",
			Command: "catch [pokemon]",
			Description: "Catch pokemon and add to your collections.",
			Call: catchCall,
		},

		"inspect": {
			Name: "Inspect",
			Command: "inspect [pokemon]",
			Description: "View details about a pokemon in your collections.",
			Call: inspectCall,
		},

		"pokedex": {
			Name: "Pokedex",
			Command: "pokedex",
			Description: "List all pokemons in your collection. (that you've caught)",
			Call: pokedexCall,
		},
	}
}