package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Item struct {
	Name        string
	Description string
	Power       int
}

type Enemy struct {
	Name   string
	Health int
	Power  int
}

type Room struct {
	Name        string
	Description string
	Items       []Item
	Enemy       *Enemy
	Exits       map[string]string
}

type Player struct {
	Health    int
	Inventory []Item
}

var rooms map[string]*Room
var currentRoom string
var player Player

func initializeGame() {
	rooms = make(map[string]*Room)

	// Initialize rooms
	rooms["entrance"] = &Room{
		Name:        "Entrance",
		Description: "You're at the entrance of a dark cave. There's a faint light coming from the north.",
		Items:       []Item{{Name: "Torch", Description: "A wooden torch", Power: 0}},
		Exits:       map[string]string{"north": "mainHall"},
	}

	rooms["mainHall"] = &Room{
		Name:        "Main Hall",
		Description: "A large hall with ancient writings on the walls. Passages lead in all directions.",
		Items:       []Item{{Name: "Sword", Description: "A rusty old sword", Power: 5}},
		Enemy:       &Enemy{Name: "Goblin", Health: 20, Power: 3},
		Exits: map[string]string{
			"north": "treasureRoom",
			"east":  "armory",
			"south": "entrance",
			"west":  "library",
		},
	}

	rooms["treasureRoom"] = &Room{
		Name:        "Treasure Room",
		Description: "A room filled with gold and jewels. But it's guarded by a dragon!",
		Enemy:       &Enemy{Name: "Dragon", Health: 100, Power: 15},
		Exits:       map[string]string{"south": "mainHall"},
	}

	rooms["armory"] = &Room{
		Name:        "Armory",
		Description: "An old armory with weapons hanging on the walls.",
		Items:       []Item{{Name: "Shield", Description: "A sturdy shield", Power: 3}},
		Exits:       map[string]string{"west": "mainHall"},
	}

	rooms["library"] = &Room{
		Name:        "Library",
		Description: "A dusty library filled with ancient books.",
		Items:       []Item{{Name: "Spellbook", Description: "A book of magic spells", Power: 7}},
		Exits:       map[string]string{"east": "mainHall"},
	}

	currentRoom = "entrance"
	player = Player{Health: 100, Inventory: []Item{}}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initializeGame()

	fmt.Println("Welcome to the Cave Adventure!")
	fmt.Println("Commands: look, go [direction], take [item], inventory, fight, quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		words := strings.Fields(input)

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "look":
			look()
		case "go":
			if len(words) > 1 {
				move(words[1])
			} else {
				fmt.Println("Go where?")
			}
		case "take":
			if len(words) > 1 {
				take(words[1])
			} else {
				fmt.Println("Take what?")
			}
		case "inventory":
			showInventory()
		case "fight":
			fight()
		case "quit":
			fmt.Println("Thanks for playing!")
			return
		default:
			fmt.Println("I don't understand that command.")
		}

		if player.Health <= 0 {
			fmt.Println("Game Over! You have died.")
			return
		}
	}
}

func look() {
	room := rooms[currentRoom]
	fmt.Println(room.Name)
	fmt.Println(room.Description)
	if len(room.Items) > 0 {
		fmt.Println("Items in the room:")
		for _, item := range room.Items {
			fmt.Printf("- %s\n", item.Name)
		}
	}
	if room.Enemy != nil {
		fmt.Printf("There's a %s here!\n", room.Enemy.Name)
	}
	fmt.Println("Exits:")
	for direction, _ := range room.Exits {
		fmt.Printf("- %s\n", direction)
	}
}

func move(direction string) {
	room := rooms[currentRoom]
	if newRoom, exists := room.Exits[direction]; exists {
		currentRoom = newRoom
		fmt.Printf("You move %s.\n", direction)
		look()
	} else {
		fmt.Printf("You can't go %s from here.\n", direction)
	}
}

func take(itemName string) {
	room := rooms[currentRoom]
	for i, item := range room.Items {
		if strings.ToLower(item.Name) == strings.ToLower(itemName) {
			player.Inventory = append(player.Inventory, item)
			room.Items = append(room.Items[:i], room.Items[i+1:]...)
			fmt.Printf("You took the %s.\n", item.Name)
			return
		}
	}
	fmt.Printf("There's no %s here.\n", itemName)
}

func showInventory() {
	if len(player.Inventory) == 0 {
		fmt.Println("Your inventory is empty.")
		return
	}
	fmt.Println("Your inventory:")
	for _, item := range player.Inventory {
		fmt.Printf("- %s (Power: %d)\n", item.Name, item.Power)
	}
	fmt.Printf("Your health: %d\n", player.Health)
}

func fight() {
	room := rooms[currentRoom]
	if room.Enemy == nil {
		fmt.Println("There's nothing to fight here.")
		return
	}

	playerPower := 1 // Base power
	for _, item := range player.Inventory {
		playerPower += item.Power
	}

	fmt.Printf("You engage in battle with the %s!\n", room.Enemy.Name)
	for {
		// Player's turn
		damage := rand.Intn(playerPower) + 1
		room.Enemy.Health -= damage
		fmt.Printf("You hit the %s for %d damage!\n", room.Enemy.Name, damage)

		if room.Enemy.Health <= 0 {
			fmt.Printf("You have defeated the %s!\n", room.Enemy.Name)
			room.Enemy = nil
			return
		}

		// Enemy's turn
		damage = rand.Intn(room.Enemy.Power) + 1
		player.Health -= damage
		fmt.Printf("The %s hits you for %d damage!\n", room.Enemy.Name, damage)

		if player.Health <= 0 {
			fmt.Println("You have been defeated!")
			return
		}

		fmt.Printf("Your health: %d, %s's health: %d\n", player.Health, room.Enemy.Name, room.Enemy.Health)
	}
}