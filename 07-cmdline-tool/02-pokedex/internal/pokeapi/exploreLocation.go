package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type ExploreLocationAreaResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"` 
		}`json:"pokemon"`        
	} `json:"pokemon_encounters"`
}


func (client *Client) ExploreLocationArea(locationName string) (ExploreLocationAreaResp, error) {
	endPoint := "/location-area/"
	url := baseURL + endPoint + locationName


	// checking cache
	cache, ok := client.Cache.Get(url);
	if ok {
		// cache hit
		var Results ExploreLocationAreaResp;
		if err := json.Unmarshal(cache, &Results); err != nil {
			return ExploreLocationAreaResp{}, err;
		}
		return Results, nil;
	}

	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return ExploreLocationAreaResp{}, err;
	}

	res, err := client.HttpClient.Do(req);
	if err != nil {
		return ExploreLocationAreaResp{}, err;
	}
	
	if res.StatusCode > 399 {
		return ExploreLocationAreaResp{}, fmt.Errorf("Operation failed. Failed to get location area data: %w", err);
	}

	defer res.Body.Close();

	bodyBytes, err := io.ReadAll(res.Body);
	if err != nil {
		return ExploreLocationAreaResp{}, err;
	}

	// cache miss
	client.Cache.Add(url, bodyBytes)

	var Results ExploreLocationAreaResp;
	if err := json.Unmarshal(bodyBytes, &Results); err != nil {
		return ExploreLocationAreaResp{}, err;
	}

	return Results, nil;
}


