package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Mrpierrouge/grimhold-go/internal/game"
)

type CLICommand struct {
	Name        string
	Description string
	Callback    func(args []string) error
}

type REPL struct {
	Game     *game.Game
	Commands map[string]CLICommand
}

var scanner = bufio.NewScanner(os.Stdin)

func (r *REPL) Run() error {
	for scanner.Scan() {
		input := strings.ToLower(scanner.Text())
		args := strings.Fields(input)
		if _, ok := r.Commands[args[0]]; !ok {
			fmt.Println("Unknown command:", input)
			continue
		}
		command := r.Commands[args[0]]
		err := command.Callback(args[1:])
		if err != nil {
			fmt.Println("Error executing command: ["+input+"] : ", err)
		}
	}
	return scanner.Err()
}

func (r *REPL) Help() {
	fmt.Println("Grimhold — a dungeon crawler")
	fmt.Println()
	fmt.Println("Commands:")

	for _, command := range r.Commands {
		fmt.Printf("  %-15s %s\n", command.Name, command.Description)
	}
}

func (r *REPL) ListCommands() []CLICommand {
	commands := []CLICommand{
		{
			Name:        "help",
			Description: "Show this help message",
			Callback: func(args []string) error {
				r.Help()
				return nil
			},
		},
		{
			Name:        "look",
			Description: "Look at the current location",
			Callback: func(s []string) error {
				r.Game.Look()
				return nil
			},
		},
		{
			Name:        "map",
			Description: "Show the map of the world",
			Callback: func(args []string) error {
				r.Game.Map()
				return nil
			},
		},
		{
			Name:        "travel",
			Description: "Travel to a new location",
			Callback: func(args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("Please provide a region ID.")
				}
				fmt.Println("Traveling to:", args[0])
				err := r.Game.Travel(args[0])
				if err != nil {
					return fmt.Errorf("Failed to travel to region: %v", err)
				}
				return nil
			},
		},
	}
	return commands
}
