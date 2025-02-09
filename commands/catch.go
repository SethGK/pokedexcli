package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

var Pokedex = make(map[string]Pokemon)

const pokemonAPIURL = "https://pokeapi.co/api/v2/pokemon"

func Catch(cache *pokecache.Cache, pokemonName string) error {
	url := fmt.Sprintf("%s/%s", pokemonAPIURL, pokemonName)
	resp, found := cache.Get(url)

	if !found {
		fmt.Println("Cache miss! Fetching from API")
		var err error
		resp, err = fetchPokemonFromAPI(url)
		if err != nil {
			return fmt.Errorf("failed to fetch pokemon data: %w", err)
		}
		cache.Add(url, resp)
	}

	var pokemonData Pokemon
	if err := json.Unmarshal(resp, &pokemonData); err != nil {
		return fmt.Errorf("failed to unmarshal pokemon data: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	catchChance := calculateCatchChance(pokemonData.BaseExperience)

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < catchChance {
		Pokedex[pokemonData.Name] = pokemonData
		fmt.Printf("%s was caught!\n", pokemonData.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
	}

	return nil
}

func fetchPokemonFromAPI(url string) ([]byte, error) {
	fmt.Println("Fetching:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching data from API: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func calculateCatchChance(baseExperience int) int {
	chance := 100 - baseExperience
	if chance < 5 {
		chance = 5
	}
	return chance
}
