package main

import (
	"fmt"

	"github.com/Mrpierrouge/grimhold-go/internal/game"
	"github.com/Mrpierrouge/grimhold-go/internal/world"
)

func main() {
	fmt.Println("Hello world")
	world, err := world.LoadWorld()
	if err != nil {
		fmt.Println("Failed to load world")
		return
	}
	fmt.Printf("Loaded world: \n")
	game := game.NewGame(world)
	game.Map()
	game.Look()
	// err = game.Travel("blackpine-forest")
	// if err != nil {
	// 	fmt.Println("Failed to travel to region2")
	// 	return
	// }
	// game.Look()
}
