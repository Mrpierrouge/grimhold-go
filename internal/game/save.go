package game

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mrpierrouge/grimhold-go/internal/world"
)

const saveFile = "save.json"

type SaveData struct {
	CurrentRegionID string         `json:"current_region_id"`
	Character       Character      `json:"character"`
	SlainBestiary   map[string]int `json:"slain_bestiary"`
}

func (g *Game) Save() error {
	data := SaveData{
		CurrentRegionID: g.CurrentRegion.ID,
		Character:       g.Character,
		SlainBestiary:   g.SlainBestiary,
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding save data: %w", err)
	}

	tmpFile := saveFile + ".tmp"

	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		return fmt.Errorf("writing temporary save: %w", err)
	}

	if err := os.Rename(tmpFile, saveFile); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("replacing save file: %w", err)
	}

	return nil
}

func (g *Game) Load() error {
	content, err := os.ReadFile(saveFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No save file found.")
			return nil
		}
		return fmt.Errorf("reading save file: %w", err)
	}

	var data SaveData
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("decoding save file: %w", err)
	}

	region, err := world.FindRegionByID(g.World, data.CurrentRegionID)
	if err != nil {
		return fmt.Errorf("loading saved region: %w", err)
	}

	g.CurrentRegion = region
	g.Character = data.Character
	g.SlainBestiary = data.SlainBestiary

	if g.SlainBestiary == nil {
		g.SlainBestiary = make(map[string]int)
	}

	fmt.Println("Game loaded successfully!")
	return nil
}
