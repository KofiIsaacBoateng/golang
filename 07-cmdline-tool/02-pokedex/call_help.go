package main

import (
	"fmt"

	"github.com/fatih/color"
)

func helpCall() {
	color.Cyan("\nWELCOME TO POKEDEX HELP!")
	bold := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	fmt.Print("Here are your available commands.\n\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf(" - %s\t%s\t%s\n",  command.Command, bold(command.Name), command.Description)
	}
	fmt.Println("")
}