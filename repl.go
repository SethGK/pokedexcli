package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/config"
	"pokedexcli/internal/models"
	"pokedexcli/internal/pokecache"
	"strings"
	"time"

	cmd "pokedexcli/commands"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config.Config) error
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
			callback:    mapCommandWrapper,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 Pokémon locations",
			callback:    mapBackCommandWrapper,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location to find Pokémon",
			callback:    commandExploreWrapper,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokémon by name",
			callback:    commandCatchWrapper,
		},
		"inspect": {
			name:        "inspect",
			description: "View details about a caught Pokémon",
			callback: func(config *config.Config) error {
				return commandInspectWrapper(config, []string{})
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Dispalys all caught Pokémon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(_ *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config.Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func mapCommandWrapper(config *config.Config) error {
	return mapCommand(config)
}

func mapBackCommandWrapper(config *config.Config) error {
	return mapBackCommand(config)
}

func commandExplore(config *config.Config, area string) error {
	fmt.Printf("Exploring %s...\n", area)

	err := cmd.Explore(config.Cache, area)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

func commandExploreWrapper(config *config.Config) error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: explore <area_name>")
		return nil
	}
	area := os.Args[1]
	return commandExplore(config, area)
}

func commandCatch(config *config.Config, pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("Usage: catch <pokemon_name>")
		return nil
	}

	err := cmd.Catch(config, config.Cache, pokemonName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func commandCatchWrapper(config *config.Config) error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: catcj <pokemon_name>")
		return nil
	}
	pokemonName := os.Args[1]
	return commandCatch(config, pokemonName)
}

func commandInspect(config *config.Config, pokemonName string) error {
	pokemon, exists := config.CaughtPokemon[pokemonName]
	if !exists {
		fmt.Println("You have not caught that Pokémon.")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Name, stat.Value)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t)
	}

	return nil
}

func commandInspectWrapper(config *config.Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: inspect <pokemon_name>")
		return nil
	}
	pokemonName := args[0]
	return commandInspect(config, pokemonName)
}

func commandPokedex(config *config.Config) error {
	if len(config.CaughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokémon first!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range config.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	return words
}

// REPL
func startRepl() {
	config := &config.Config{
		Cache:         pokecache.NewCache(10 * time.Second),
		CaughtPokemon: make(map[string]models.Pokemon),
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
			} else if commandName == "catch" {
				if len(args) == 0 {
					fmt.Println("Usage: catch <pokemon_name>")
					continue
				}
				err := commandCatch(config, args[0])
				if err != nil {
					fmt.Println("Error:", err)
				}
				continue
			} else if commandName == "inspect" {
				if len(args) == 0 {
					fmt.Println("Usage: inspect <pokemon_name>")
					continue
				}
				err := commandInspectWrapper(config, args)
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
