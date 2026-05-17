package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationAreaResp struct {
	Count    int64   `json:"count"`
	Next     *string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (client *Client) LocationAreas(uri *string) (LocationAreaResp, error) {
	endPoint := "/location-area"
	url := baseURL + endPoint

	if uri != nil {
		url = *uri;
	}


	// checking cache
	cache, ok := client.Cache.Get(url);
	if ok {
		// cache hit
		var Results LocationAreaResp;
		if err := json.Unmarshal(cache, &Results); err != nil {
			return LocationAreaResp{}, err;
		}
		return Results, nil;
	}

	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return LocationAreaResp{}, err;
	}

	res, err := client.HttpClient.Do(req);
	if err != nil {
		return LocationAreaResp{}, err;
	}
	
	if res.StatusCode > 399 {
		return LocationAreaResp{}, fmt.Errorf("Operation failed. Failed to fetch location areas: %w", err);
	}

	defer res.Body.Close();

	bodyBytes, err := io.ReadAll(res.Body);
	if err != nil {
		return LocationAreaResp{}, err;
	}
	// cache miss
	client.Cache.Add(url, bodyBytes)

	var Results LocationAreaResp;
	if err := json.Unmarshal(bodyBytes, &Results); err != nil {
		return LocationAreaResp{}, err;
	}

	return Results, nil;
}