package restful_api_weather

import "context"

type Weather struct {
	City        string
	Temperature float64
	Wind        float64
}

type WeatherProvider interface {
	FetchWeatherReport(ctx context.Context, city string) (*Weather, error)
}

type WeatherPublisher interface {
	PublishWeatherReport(ctx context.Context, weather *Weather) error
}

type WeatherService struct {
	api WeatherProvider
	pub WeatherPublisher
}

func NewWeatherService(api WeatherProvider, pub WeatherPublisher) *WeatherService {
	return &WeatherService{
		api: api,
		pub: pub,
	}
}

func (s *WeatherService) GetWeatherReport(ctx context.Context, city string) (*Weather, error) {
	return s.api.FetchWeatherReport(ctx, city)
}
