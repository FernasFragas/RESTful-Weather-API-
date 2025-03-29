package weatherservice

import (
	"context"
	"fmt"
	"strings"
)

type ReporterProvider[T any] interface {
	FetchReportData(ctx context.Context, city string) (*DataToReport[T], error)
	FetchGeneralInfo(ctx context.Context, city string) (*DataToReport[T], error)
}

type ReporterPublisher[T any] interface {
	PublishReportData(ctx context.Context, weather *GeneralWeatherInfo) error
}

type DataToReport[T any] struct {
	Data T
}

type GeneralWeatherInfo struct {
	Weather
	Waves
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	EmbedURL string  `json:"embed_url"`
}

type WeatherReporters struct {
	weatherApi ReporterProvider[GeneralWeatherInfo]
	wavesApi   ReporterProvider[GeneralWeatherInfo]
}

func NewWeatherReporters(weatherApi ReporterProvider[GeneralWeatherInfo], wavesApi ReporterProvider[GeneralWeatherInfo]) *WeatherReporters {
	return &WeatherReporters{
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

func (s *WeatherReporters) GenerateReport(ctx context.Context, city string) (*GeneralWeatherInfo, error) {
	weatherInfo, err := s.weatherApi.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	cityCoordinates := fmt.Sprintf("%f,%f", weatherInfo.Data.Lat, weatherInfo.Data.Lon)

	waveInfo, err := s.wavesApi.FetchReportData(ctx, cityCoordinates)
	if err != nil {
		return nil, err
	}

	return &GeneralWeatherInfo{
		City:     weatherInfo.Data.City,
		Country:  strings.ToLower(weatherInfo.Data.Country),
		Lat:      weatherInfo.Data.Lat,
		Lon:      weatherInfo.Data.Lon,
		Waves:    waveInfo.Data.Waves,
		Weather:  weatherInfo.Data.Weather,
		EmbedURL: fmt.Sprintf("https://embed.windy.com/embed2.html?lat=%f&lon=%f&zoom=11&level=surface&overlay=wind", weatherInfo.Data.Lat, weatherInfo.Data.Lon),
	}, nil

}

func NewVideoStreamReporters(youtubeApi ReporterProvider[VideosStream]) *VideoStreamReporters {
	return &VideoStreamReporters{
		youtubeApi: youtubeApi,
	}
}

type VideoStreamReporters struct {
	youtubeApi ReporterProvider[VideosStream]
}

type GeneralVideoStreamInfo struct {
	Title   string
	VideoID string
}

type VideosStream []GeneralVideoStreamInfo

func (s *VideoStreamReporters) GenerateReport(ctx context.Context, city string) (*VideosStream, error) {
	videos, err := s.youtubeApi.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}
	return &videos.Data, nil
}

type Hotels []Hotel

type HotelsApi struct {
	hotelsApi ReporterProvider[Hotels]
	photosAPI ReporterProvider[HotelsPhotos]
}

func NewHotelsApi(hotelsApi ReporterProvider[Hotels]) *HotelsApi {
	return &HotelsApi{
		hotelsApi: hotelsApi,
	}
}

type Hotel struct {
	HotelName   string
	HotelURL    string
	HotelPhotos []string
}

func (s *HotelsApi) GenerateReport(ctx context.Context, city string) (*Hotels, error) {
	hotels, err := s.hotelsApi.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	// search for the photos
	hotelsPhotos, err := s.photosAPI.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	for _, hotel := range hotels.Data {
		for _, hotelPhoto := range hotelsPhotos.Data {
			if hotel.HotelName == hotelPhoto.HotelName {
				hotel.HotelPhotos = hotelPhoto.HotelPhotos
			}
		}
	}

	return &hotels.Data, nil
}

type HotelsPhotos []HotelPhotos

type HotelPhotos struct {
	HotelName   string
	HotelPhotos []string
}
