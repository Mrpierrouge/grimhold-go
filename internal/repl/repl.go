package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	Scanner  *bufio.Scanner
}

func (r *REPL) ReadInput() (string, error) {
	if !r.Scanner.Scan() {
		return "", r.Scanner.Err()
	}

	return strings.ToLower(r.Scanner.Text()), nil
}
func (r *REPL) Run() error {
	if r.Scanner == nil {
		r.Scanner = bufio.NewScanner(os.Stdin)
	}
	for {
		input, err := r.ReadInput()
		if err != nil {
			return err
		}
		args := strings.Fields(input)
		if _, ok := r.Commands[args[0]]; !ok {
			fmt.Println("Unknown command:", input)
			continue
		}
		command := r.Commands[args[0]]
		errCallback := command.Callback(args[1:])
		if errCallback != nil {
			fmt.Println("Error executing command: ["+input+"] : ", errCallback)
		}
	}
}

func (r *REPL) Help() {
	fmt.Println("Grimhold — a dungeon crawler")
	fmt.Println()
	fmt.Println("Commands:")

	for _, command := range r.Commands {
		fmt.Printf("  %-15s %s\n", command.Name, command.Description)
	}
}

func (r *REPL) Bestiary() error {
	page, err := r.Game.Client.GetMonsters("")
	if err != nil {
		return fmt.Errorf("loading bestiary: %w", err)
	}

	for {
		fmt.Println("\n--- Grimhold Bestiary ---")

		for _, monster := range page.Results {
			fmt.Printf("Name: %s \n", monster.Name)
		}

		fmt.Println("\nNavigation:")

		if page.Previous != nil {
			fmt.Print("[prev] ")
		}

		if page.Next != nil {
			fmt.Print("[next] ")
		}

		fmt.Println("[stop]")
		fmt.Print("bestiary > ")

		input, err := r.ReadInput()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch strings.ToLower(input) {
		case "next":
			if page.Next == nil {
				fmt.Println("You are already on the last page.")
				continue
			}
			page, err = r.Game.Client.GetMonsters(*page.Next)
			if err != nil {
				return fmt.Errorf("fetching next page: %w", err)
			}

		case "prev":
			if page.Previous == nil {
				fmt.Println("You are already on the first page.")
				continue
			}

			page, err = r.Game.Client.GetMonsters(*page.Previous)
			if err != nil {
				return fmt.Errorf("fetching previous page: %w", err)
			}

		case "stop":
			fmt.Println("Leaving the bestiary...")
			return nil

		default:
			fmt.Println("Unknown option. Use: prev, next or stop.")
			continue
		}

		if err != nil {
			return fmt.Errorf("changing bestiary page: %w", err)
		}
	}
}

func (r *REPL) Fight() error {
	if r.Game.CurrentFoe == nil {
		return fmt.Errorf("there is no active enemy; use explore first")
	}

	fmt.Println("\nCombat started!")
	fmt.Println("Commands: fight | flee")

	for r.Game.CurrentFoe != nil {
		fmt.Print("\ncombat > ")

		input, err := r.ReadInput()
		if err != nil {
			return err
		}

		switch strings.ToLower(input) {
		case "fight":
			defeated := r.Game.PlayerAttack()

			if defeated {
				r.Game.Victory()
				return nil
			}

			r.Game.FoeAttack()

			if r.Game.Character.HP <= 0 {
				r.Game.Defeat()
				return nil
			}

		case "flee":
			escaped := r.Game.Flee()

			if escaped {
				return nil
			}

			if r.Game.Character.HP <= 0 {
				r.Game.Defeat()
				return nil
			}

		default:
			fmt.Println("Use fight or flee.")
		}
	}

	return nil
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
				err := r.Game.Travel(args[0])
				if err != nil {
					return fmt.Errorf("Failed to travel to region: %v", err)
				}
				return nil
			},
		},
		{
			Name:        "bestiary",
			Description: "Show a list of monsters from the Open5e API",
			Callback: func(args []string) error {
				return r.Bestiary()
			},
		},
		{
			Name:        "inspect",
			Description: "Inspect a monster by name",
			Callback: func(args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("Please provide a monster name.")
				}
				monsterName := strings.NewReplacer(
					" ", "+",
					"-", "+",
					"_", "+",
				).Replace(strings.Join(args, " "))
				monster, err := r.Game.Client.GetMonster(monsterName)
				if err != nil {
					return fmt.Errorf("Failed to inspect monster: %v", err)
				}
				fmt.Printf("Name: %s\n", monster.Name)
				fmt.Printf("HP: %d\n", monster.HP)
				fmt.Printf("Armor: %d\n", monster.AC)
				fmt.Printf("XP: %d\n", monster.Experience)
				return nil
			},
		},
		{
			Name:        "character",
			Description: "Show your character sheet",
			Callback: func(args []string) error {
				r.Game.Character.Show()
				return nil
			},
		},
		{
			Name:        "rest",
			Description: "Restore your HP",
			Callback: func(args []string) error {
				return r.Game.Character.Rest()
			},
		},
		{
			Name:        "explore",
			Description: "Explore the current region and search for enemies",
			Callback: func(args []string) error {
				return r.Game.Explore()
			},
		},
		{
			Name:        "fight",
			Description: "Fight the current enemy",
			Callback: func(args []string) error {
				return r.Fight()
			},
		},
	}
	return commands
}
