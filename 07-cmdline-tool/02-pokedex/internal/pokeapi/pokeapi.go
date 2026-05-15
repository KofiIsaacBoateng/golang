package pokeapi

import (
	"net/http"
	"pokedex/internal/pokecache"
	"time"
)


type Client struct {
	Cache pokecache.Cache
	HttpClient *http.Client
}

const baseURL = "https://pokeapi.co/api/v2"


func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		HttpClient: &http.Client{},
		Cache: pokecache.NewCache(cacheInterval),
	}
}

