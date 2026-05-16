package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`

		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct{
				TimeEpoch int `json:"time_epoch"`
				TempC float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Accra";
	if len(os.Args) >= 2 {
		q = strings.Join(os.Args[1:], "%20");
	}
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=d94208522b5a47e4b99122520261905&q=" + q + "&days=1&aqi=no&alerts=no");
	if err != nil {
		log.Fatalf("Unable to get weather data, please try again!: %v", err);
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Unable to get weather data, please try again!");
	}

	defer res.Body.Close();

	bodyBytes, err := io.ReadAll(res.Body);
	if(err != nil ){
		log.Fatalf("Failed to parse body bytes: %v", err);
	}

	var weather Weather;
	if err := json.Unmarshal(bodyBytes, &weather); err != nil {
		log.Fatalf("It's not you! It is us. Please try again! %v", err);
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour;

	fmt.Printf(
		"%s, %s: %.0fC, %s\n", 
		location.Name, 
		location.Country, 
		current.TempC, 
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0);

		if date.Before(time.Now()) {
			continue;
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if(hour.ChanceOfRain < 40) {
			fmt.Print(message)
		}else {
			color.Red(message)
		}
	}
}