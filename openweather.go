// Package openweather package provides a simple API to fetch weather reports for a given city.
package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"weatherservice"
)

const openWeatherMapWebhookURL = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherAPI struct {
	client *http.Client

	key string
}

func NewWeatherAPI(key string) *WeatherAPI {
	return &WeatherAPI{
		client: http.DefaultClient,
		key:    key,
	}
}

func (api *WeatherAPI) FetchWeatherReport(ctx context.Context, city string) (*weatherservice.Weather, error) {
	if api.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	queryParams := url.Values{}
	queryParams.Add("q", city)
	queryParams.Add("units", "metric")
	queryParams.Add("APPID", api.key)

	apiUrl := fmt.Sprintf("%s%s", openWeatherMapWebhookURL, queryParams.Encode())

	resp, err := api.client.Get(apiUrl)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var weather weatherData

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weatherservice.Weather{
		City:        weather.Name,
		Temperature: weather.Temperature.Temperature,
		FeelsLike:   weather.Temperature.FeelsLike,
		Wind:        weather.Wind.Speed,
		Humidity:    weather.Temperature.Humidity,
		Condition:   weather.Weather[0].Description,
		Country:     strings.ToLower(weather.Sys.Country),
	}, nil
}

// unexported
type weatherData struct {
	Weather     weatherInfo `json:"weather"`
	Temperature temperature `json:"main"`
	Wind        wind        `json:"wind"`
	Sys         sys         `json:"sys"`
	Name        string      `json:"name"`
}

type weatherInfo []struct {
	WeatherState string `json:"main"`
	Description  string `json:"description"`
}

type temperature struct {
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"tem_max"`
	Humidity    float64 `json:"humidity"`
}

type sys struct {
	Country string  `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset  float64 `json:"sunset"`
}

type wind struct {
	Speed float64 `json:"speed"`
}
