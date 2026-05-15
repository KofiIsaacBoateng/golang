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
	Call func()
}


func Repl() {
	for {
		scanner := bufio.NewScanner(os.Stdin);

		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgHiGreen).SprintFunc()

		fmt.Printf("%s %s ", cyan("pokedex"), green("$>"))

		scanner.Scan()
		input := strings.TrimSpace(scanner.Text());

		// onPressEnter
		if input == "" {
			continue;
		}

		commands := getCommands();

		command, ok := commands[input];
		if !ok {
			// invalid input
			fmt.Printf("command: %s is INVALID.\n", input)
			continue;
		}

		command.Call();
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
	}
}