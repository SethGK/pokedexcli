package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

const pokeAPIURL = "https://pokeapi.co/api/v2/location-area/"

func Explore(cache *pokecache.Cache, area string) error {
	url := fmt.Sprintf("%s%s", pokeAPIURL, area)

	resp, found := cache.Get(url)
	if !found {
		fmt.Println("Cache miss! Fetching from API")
		var err error
		resp, err = fetchFromAPI(url)
		if err != nil {
			return fmt.Errorf("failed to fetch location data: %w", err)
		}
		cache.Add(url, resp)
	}

	var locationData LocationAreaResponse
	if err := json.Unmarshal(resp, &locationData); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(locationData.PokemonEncounters) == 0 {
		return fmt.Errorf("no Pokémon found in location area: %s", area)
	}

	fmt.Println("Exploring", area, "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationData.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func fetchFromAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching data from API: %s", resp.Status)
	}

	var locationData LocationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&locationData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Raw response body:", string(body))

	return json.Marshal(locationData)
}

type LocationAreaResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
				Chance   int `json:"chance"`
				Method   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
