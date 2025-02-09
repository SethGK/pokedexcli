package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokedexcli/internal/config"
	"pokedexcli/internal/pokecache"
	"time"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

type APIResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

var cache = pokecache.NewCache(10 * time.Second)

func fetchLocations(url string) (*APIResponse, error) {
	if data, found := cache.Get(url); found {
		fmt.Println("Cache hit!")
		var cachedResponse APIResponse

		if err := json.Unmarshal(data, &cachedResponse); err != nil {
			return nil, err
		}
		return &cachedResponse, nil
	}
	fmt.Println("Cache miss! Fetching from API")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err == nil {
		cache.Add(url, jsonData)
	}

	return &data, nil
}

func mapCommand(config *config.Config) error { // Change to *config.Config
	data, err := fetchLocations(getOrDefault(config.Next, baseURL))
	if err != nil {
		return fmt.Errorf("error fetching locations: %w", err)
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous
	return nil
}

func mapBackCommand(config *config.Config) error { // Change to *config.Config
	if config.Previous == nil {
		fmt.Println("No previous page")
		return nil
	}

	data, err := fetchLocations(*config.Previous)
	if err != nil {
		return fmt.Errorf("error fetching locations: %w", err)
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous
	return nil
}

func getOrDefault(url *string, defaultURL string) string {
	if url == nil {
		return defaultURL
	}
	return *url
}
