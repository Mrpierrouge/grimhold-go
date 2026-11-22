package main

import (
	"fmt"

	"github.com/Mrpierrouge/grimhold-go/internal/game"
	"github.com/Mrpierrouge/grimhold-go/internal/open5e"
	"github.com/Mrpierrouge/grimhold-go/internal/repl"
	"github.com/Mrpierrouge/grimhold-go/internal/world"
)

func main() {
	world, err := world.LoadWorld()
	if err != nil {
		fmt.Println("Failed to load world")
		return
	}
	fmt.Printf("Loaded world: \n")
	game := game.NewGame(world)
	if game == nil {
		fmt.Println("Failed to create new game")
		return
	}
	client := open5e.NewClient()
	defer client.Close()

	repl := repl.REPL{
		Game:     game,
		Commands: make(map[string]repl.CLICommand),
	}
	var commands = repl.ListCommands()

	for _, cmd := range commands {
		repl.Commands[cmd.Name] = cmd
	}
	err = repl.Run()
}
