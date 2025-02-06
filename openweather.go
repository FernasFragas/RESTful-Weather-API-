package restful_api_weather

//This comunicates with that specific api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const openWeatherMapWebhookURL = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherAPI struct {
	client *http.Client

	key string
}

func NewWeatherAPI(key string) *WeatherAPI {
	return &WeatherAPI{
		key: key,
	}
}

func (api *WeatherAPI) FetchWeatherReport(ctx context.Context, city string) (*Weather, error) {
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

	var weather WeatherData

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &Weather{
		City:        weather.Name,
		Temperature: weather.Temperature.TempMax,
		Wind:        0,
	}, nil
}

// unexported
type WeatherData struct {
	Weather     WeatherInfo `json:"weather"`
	Temperature Temperature `json:"main"`
	Wind        Wind        `json:"wind"`
	Sys         Sys         `json:"sys"`
	Name        string      `json:"name"`
}

type WeatherInfo []struct {
	WeatherState string `json:"main"`
	Description  string `json:"description"`
}

type Temperature struct {
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"tem_max"`
	Humidity  float64 `json:"humidity"`
}

type Sys struct {
	Country string  `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset  float64 `json:"sunset"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}
