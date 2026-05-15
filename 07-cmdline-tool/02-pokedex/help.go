package main

import (
	"fmt"

	"github.com/fatih/color"
)

func helpCall(cfg *Config) {
	color.Cyan("\nWELCOME TO POKEDEX HELP!")
	fmt.Print("Here are your available commands.\n\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf(" - %s\t%s\t%s\n",  command.Command, command.Name, command.Description)
	}
	fmt.Println("")
}