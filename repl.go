package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"strings"
	"time"

	cmd "pokedexcli/commands"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 Pokémon locations",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 Pokémon locations",
			callback:    mapBackCommand,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location to find Pokémon",
			callback:    commandExploreWrapper,
		},
	}
}

func commandExit(_ *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExplore(config *Config, area string) error {
	fmt.Printf("Exploring %s...\n", area)

	err := cmd.Explore(config.Cache, area)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

func commandExploreWrapper(config *Config) error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: explore <area_name>")
		return nil
	}
	area := os.Args[1]
	return commandExplore(config, area)
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	return words
}

// REPL
func startRepl() {
	config := &Config{
		Cache: pokecache.NewCache(10 * time.Second),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		commandName := input[0]
		args := input[1:]

		if cmd, exists := commands[commandName]; exists {
			if commandName == "explore" {
				if len(args) == 0 {
					fmt.Println("Usage: explore <area_name>")
					continue
				}
				err := commandExplore(config, args[0])
				if err != nil {
					fmt.Println("Error:", err)
				}
				continue
			}
			err := cmd.callback(config)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
