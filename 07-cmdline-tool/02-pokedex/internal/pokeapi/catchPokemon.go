package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type PokemonResp struct {
	ID                     int64         `json:"id"`                      
	Name                   string        `json:"name"`                    
	BaseExperience         int64         `json:"base_experience"`         
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Height                 int64         `json:"height"`                  
	Weight                 int64         `json:"weight"`                  
	Stats                  []struct {
								BaseStat int64   `json:"base_stat"`
								Effort   int64   `json:"effort"`   
								Stat     struct {
									Name string `json:"name"`
									URL  string `json:"url"` 
								} `json:"stat"`     
							} `json:"stats"`                   

}


func (client *Client) CatchPokemon(pokemonName string) (PokemonResp, error) {
	endPoint := "/pokemon/"
	url := baseURL + endPoint + pokemonName


	// checking cache
	cache, ok := client.Cache.Get(url);
	if ok {
		// cache hit
		var Results PokemonResp;
		if err := json.Unmarshal(cache, &Results); err != nil {
			return PokemonResp{}, err;
		}
		return Results, nil;
	}

	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return PokemonResp{}, err;
	}

	res, err := client.HttpClient.Do(req);
	if err != nil {
		return PokemonResp{}, err;
	}
	
	if res.StatusCode > 399 {
		return PokemonResp{}, fmt.Errorf("Operation failed. Failed to get pokemon data: %w", err);
	}

	defer res.Body.Close();

	bodyBytes, err := io.ReadAll(res.Body);
	if err != nil {
		return PokemonResp{}, err;
	}
	// cache miss
	client.Cache.Add(url, bodyBytes)

	var Results PokemonResp;
	if err := json.Unmarshal(bodyBytes, &Results); err != nil {
		return PokemonResp{}, err;
	}

	return Results, nil;
}


