package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Next     *string
	Previous *string
}

const baseURL = "https://pokeapi.co/api/v2/location-area/"

type APIResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

func fetchLocations(url string) (*APIResponse, error) {
	if url == "" {
		url = baseURL
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func mapCommand(config *Config) error {
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

func mapBackCommand(config *Config) error {
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
