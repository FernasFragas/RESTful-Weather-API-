package weatherservice

import (
	"context"
	"fmt"
	"strings"
)

type ReporterProvider[T any] interface {
	FetchReportData(ctx context.Context, city string) (*GeneralInfo, error)
	FetchGeneralInfo(ctx context.Context, city string) (*GeneralInfo, error)
}

type ReporterPublisher[T any] interface {
	PublishReportData(ctx context.Context, weather *GeneralInfo) error
}

type GeneralInfo struct {
	Weather
	Waves
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	EmbedURL string  `json:"embed_url"`
}

type Reporters struct {
	weatherApi ReporterProvider[Weather]
	wavesApi   ReporterProvider[Waves]
}

func NewReporters(weatherApi ReporterProvider[Weather], wavesApi ReporterProvider[Waves]) *Reporters {
	return &Reporters{
		weatherApi: weatherApi,
		wavesApi:   wavesApi,
	}
}

type Weather struct {
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	Wind        float64 `json:"wind"`
	Humidity    float64 `json:"humidity"`
	Condition   string  `json:"condition"`
}

type Waves struct {
	Height float64 `json:"height"`
}

func (s *Reporters) GenerateReport(ctx context.Context, city string) (*GeneralInfo, error) {
	weatherInfo, err := s.weatherApi.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	cityCoordinates := fmt.Sprintf("%f,%f", weatherInfo.Lat, weatherInfo.Lon)

	waveInfo, err := s.wavesApi.FetchReportData(ctx, cityCoordinates)
	if err != nil {
		return nil, err
	}

	return &GeneralInfo{
		City:     weatherInfo.City,
		Country:  strings.ToLower(weatherInfo.Country),
		Lat:      weatherInfo.Lat,
		Lon:      weatherInfo.Lon,
		Waves:    waveInfo.Waves,
		Weather:  weatherInfo.Weather,
		EmbedURL: fmt.Sprintf("https://embed.windy.com/embed2.html?lat=%f&lon=%f&zoom=5&level=surface&overlay=wind", weatherInfo.Lat, weatherInfo.Lon),
	}, nil

}
