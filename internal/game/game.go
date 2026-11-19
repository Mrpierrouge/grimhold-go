package game

import (
	"errors"
	"fmt"

	"github.com/Mrpierrouge/grimhold-go/internal/world"
)

type Game struct {
	World         world.World
	CurrentRegion world.Region
}

func NewGame(w world.World) *Game {
	startRegion, err := world.FindRegionByID(w, w.Start)
	if err != nil {
		return nil
	}
	return &Game{
		World:         w,
		CurrentRegion: startRegion,
	}
}

func (g *Game) Map() {
	fmt.Printf("Known regions :\n\n")
	for _, region := range g.World.Regions {
		if region.ID == g.CurrentRegion.ID {
			fmt.Printf("Region: %s (current)\n", region.Name)
		} else {
			fmt.Printf("Region: %s\n", region.Name)
		}
		fmt.Printf("Danger: %d\n", region.Danger)
		fmt.Printf("Foes: %v\n", region.Foes)
		fmt.Printf("Exits: %v\n\n", region.Exits)
	}
}

func (g *Game) Look() error {
	fmt.Println("You Are in : " + g.CurrentRegion.Name + "\n")
	fmt.Println(g.CurrentRegion.Description + "\n\n")
	fmt.Println("Danger Level: " + fmt.Sprint(g.CurrentRegion.Danger) + "\n")
	fmt.Println(fmt.Sprint(g.CurrentRegion.Foes) + "\n\n")
	fmt.Println("Exits: " + fmt.Sprint(g.CurrentRegion.Exits) + "\n")
	return nil
}

func (g *Game) Travel(regionID string) error {
	for _, exit := range g.CurrentRegion.Exits {
		if exit == regionID {
			newRegion, err := world.FindRegionByID(g.World, regionID)
			if err != nil {
				return err
			}
			g.CurrentRegion = newRegion
			return nil
		}
	}

	return errors.New("You can't travel to that region from here.")

}
