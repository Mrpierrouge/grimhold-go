package game

import (
	"fmt"
	"math/rand"

	"github.com/Mrpierrouge/grimhold-go/internal/open5e"
	"github.com/Mrpierrouge/grimhold-go/internal/world"
)

func rollDice(sides int) int {
	return rand.Intn(sides) + 1
}

func rollD20() int {
	return rollDice(20)
}

func rollDamage() int {
	return rollDice(6)
}

func (g *Game) StartCombat(foe open5e.Monster) {
	g.CurrentFoe = &foe
	fmt.Printf("\nA %s appears!\n", foe.Name)
}

func (g *Game) Explore() error {
	if len(g.CurrentRegion.Foes) == 0 {
		return fmt.Errorf("there are no foes in this region")
	}

	rolls := 1
	if g.CurrentRegion.Danger >= 3 {
		rolls = 2
	}
	if g.CurrentRegion.Danger >= 5 {
		rolls = 3
	}

	var selected *open5e.Monster

	for i := 0; i < rolls; i++ {
		name := g.CurrentRegion.Foes[rand.Intn(len(g.CurrentRegion.Foes))]

		foe, err := g.Client.GetMonster(name)
		if err != nil {
			return fmt.Errorf("fetching %s: %w", name, err)
		}

		if selected == nil || foe.ChallengeRating > selected.ChallengeRating {
			copy := foe
			selected = &copy
		}
	}

	g.StartCombat(*selected)
	fmt.Println("Type 'fight' to fight or 'flee' to escape.")
	return nil
}
func (g *Game) PlayerAttack() bool {
	foe := g.CurrentFoe
	if foe == nil {
		return false
	}

	roll := rollD20()
	total := roll + g.Character.AttackBonus

	fmt.Printf("You roll %d + %d = %d\n",
		roll, g.Character.AttackBonus, total)

	if total < foe.AC {
		fmt.Println("Your attack misses!")
		return false
	}

	damage := rollDamage()
	foe.HP -= damage

	if foe.HP < 0 {
		foe.HP = 0
	}

	fmt.Printf("You hit %s for %d damage! (%d HP remaining)\n",
		foe.Name, damage, foe.HP)

	return foe.HP <= 0
}
func (g *Game) FoeAttack() {
	foe := g.CurrentFoe
	if foe == nil {
		return
	}

	roll := rollD20()
	total := roll + 2

	playerAC := g.Character.Armor

	fmt.Printf("%s rolls %d + 2 = %d\n", foe.Name, roll, total)

	if total < playerAC {
		fmt.Printf("%s misses you!\n", foe.Name)
		return
	}

	damage := rollDice(4)
	g.Character.HP -= damage

	if g.Character.HP < 0 {
		g.Character.HP = 0
	}

	fmt.Printf("%s hits you for %d damage! (%d HP remaining)\n",
		foe.Name, damage, g.Character.HP)
}
func (g *Game) Flee() bool {
	roll := rollD20() + g.Character.Dexterity

	fmt.Printf("You try to flee! Roll: %d\n", roll)

	if roll >= 12 {
		fmt.Println("You successfully escaped!")
		g.CurrentFoe = nil
		return true
	}

	fmt.Println("You failed to escape!")
	g.FoeAttack()
	return false
}
func (g *Game) Defeat() {
	fmt.Println("\nYou have been defeated...")
	fmt.Println("You wake up at the Crossroads with 1 HP.")

	g.Character.HP = 1
	g.CurrentFoe = nil

	region, err := world.FindRegionByID(g.World, "crossroads")
	if err == nil {
		g.CurrentRegion = region
	}
}
func (g *Game) Victory() {
	foe := g.CurrentFoe
	if foe == nil {
		return
	}

	fmt.Printf("\nYou defeated %s!\n", foe.Name)
	fmt.Printf("You gain %d XP.\n", foe.Experience)
	g.Character.AddXP(foe.Experience)
	g.SlainBestiary[foe.Name]++
	g.CurrentFoe = nil
}
