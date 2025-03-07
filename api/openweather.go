// Package weatherservice package provides a simple API to fetch weather reports for a given city.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (api *WeatherAPI) FetchReportData(ctx context.Context, city string) (*weatherservice.GeneralInfo, error) {
	weather, err := api.fetchLocationInfo(ctx, city)
	if err != nil {
		return nil, err
	}

	var condition string
	if len(weather.Weather) > 0 {
		condition = weather.Weather[0].Description
	}

	return &weatherservice.GeneralInfo{
		City:    weather.Name,
		Country: weather.Sys.Country,
		Lon:     weather.Coordinates.Lon,
		Lat:     weather.Coordinates.Las,
		Weather: weatherservice.Weather{
			Temperature: weather.Temperature.Temperature,
			FeelsLike:   weather.Temperature.FeelsLike,
			Wind:        weather.Wind.Speed,
			Humidity:    weather.Temperature.Humidity,
			Condition:   condition,
		},
	}, nil
}

func (api *WeatherAPI) FetchGeneralInfo(ctx context.Context, city string) (*weatherservice.GeneralInfo, error) {
	generalInfo, err := api.fetchLocationInfo(ctx, city)
	if err != nil {
		return nil, err
	}

	return &weatherservice.GeneralInfo{
		City:    generalInfo.Name,
		Country: generalInfo.Sys.Country,
		Lon:     generalInfo.Coordinates.Lon,
		Lat:     generalInfo.Coordinates.Las,
	}, nil
}

func (api *WeatherAPI) fetchLocationInfo(ctx context.Context, city string) (*weatherData, error) {
	if api.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	apiUrl, err := api.setupQueryParams(city)
	if err != nil {
		return nil, err
	}

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

	return &weather, nil
}

func (api *WeatherAPI) setupQueryParams(city string) (string, error) {
	if api.client == nil {
		return "", fmt.Errorf("client not initialized")
	}

	queryParams := url.Values{}
	queryParams.Add("q", city)
	queryParams.Add("units", "metric")
	queryParams.Add("APPID", api.key)

	apiUrl := fmt.Sprintf("%s%s", openWeatherMapWebhookURL, queryParams.Encode())

	return apiUrl, nil
}

// unexported
type weatherData struct {
	Coordinates Coordinates `json:"coord"`
	Weather     weatherInfo `json:"weather"`
	Temperature temperature `json:"main"`
	Wind        wind        `json:"wind"`
	Sys         sys         `json:"sys"`
	Name        string      `json:"name"`
}

type Coordinates struct {
	Lon float64 `json:"lon"`
	Las float64 `json:"lat"`
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
