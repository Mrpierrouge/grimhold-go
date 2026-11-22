package game

import "fmt"

type Character struct {
	HP          int
	MaxHP       int
	Level       int
	XP          int
	AttackBonus int
	Weapon      string

	RestedInRegion bool
}

func NewCharacter() Character {
	return Character{
		HP:          20,
		MaxHP:       20,
		Level:       1,
		XP:          0,
		AttackBonus: 2,
		Weapon:      "Rusty Sword",
	}
}

func (c *Character) Show() {
	fmt.Println("\n===== YOU =====")
	fmt.Printf("Level:        %d\n", c.Level)
	fmt.Printf("XP:           %d / %d\n", c.XP, c.Level*c.Level*50)
	fmt.Printf("HP:           %d / %d\n", c.HP, c.MaxHP)
	fmt.Printf("Attack Bonus: +%d\n", c.AttackBonus)
	fmt.Printf("Weapon:       %s\n", c.Weapon)
}
func (c *Character) Rest() error {
	if c.RestedInRegion {
		return fmt.Errorf("you have already rested in this region")
	}

	if c.HP == c.MaxHP {
		return fmt.Errorf("you already have full health")
	}

	c.HP = c.MaxHP
	c.RestedInRegion = true

	fmt.Println("You rest and recover all your HP.")
	return nil
}

func (c *Character) AddXP(amount int) {
	c.XP += amount

	fmt.Printf("You gained %d XP!\n", amount)

	for c.XP >= c.Level*c.Level*50 {
		c.Level++
		c.MaxHP += 5
		c.HP += 5
		c.AttackBonus++

		fmt.Printf("\nLEVEL UP! You are now level %d!\n", c.Level)
		fmt.Println("Max HP increased by 5.")
		fmt.Println("Attack bonus increased by 1.")
	}
}
