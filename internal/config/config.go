package config

import (
	"pokedexcli/internal/models"
	"pokedexcli/internal/pokecache"
)

type Config struct {
	Cache         *pokecache.Cache
	CaughtPokemon map[string]models.Pokemon
	Next          *string
	Previous      *string
}
