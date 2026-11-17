package world

import (
	"encoding/json"
	"errors"
	"os"
)

type Region struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Danger      int      `json:"danger"`
	Foes        []string `json:"foes"`
	Exits       []string `json:"exits"`
}

type World struct {
	Start   string   `json:"start"`
	Regions []Region `json:"regions"`
}

func FindRegionByID(world World, id string) (Region, error) {
	for _, region := range world.Regions {
		if region.ID == id {
			return region, nil
		}
	}
	return Region{}, errors.New("Region not found")
}

func LoadWorld() (World, error) {
	data, err := os.ReadFile("internal/world/world.json")
	if err != nil {
		return World{}, errors.New("Error reading world.json: " + err.Error())
	}

	var world World
	if err := json.Unmarshal(data, &world); err != nil {
		return World{}, errors.New("Error unmarshalling world.json: " + err.Error())
	}

	return world, nil
}
