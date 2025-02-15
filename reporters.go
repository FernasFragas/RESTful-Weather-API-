package weatherservice

import "context"

type InfoDataProvider[T any] struct {
	Data T
}

type InfoProvider[T any] interface {
	FetchWeatherReport(ctx context.Context, city string) (*InfoDataProvider[T], error)
}

type InfoPublisher[T any] interface {
	PublishWeatherReport(ctx context.Context, weather *InfoDataProvider[T]) error
}

type Weather struct {
	City        string
	Temperature float64
	FeelsLike   float64
	Wind        float64
	Humidity    float64
	Condition   string
	Country     string
}

type WeatherService struct {
	api InfoProvider[Weather]
	pub InfoPublisher[Weather]
}

func NewWeatherService(api InfoProvider[Weather], pub InfoPublisher[Weather]) *WeatherService {
	return &WeatherService{
		api: api,
		pub: pub,
	}
}

func (s *WeatherService) GetWeatherReport(ctx context.Context, city string) (*InfoDataProvider[Weather], error) {
	return s.api.FetchWeatherReport(ctx, city)
}

type Waves struct {
	Height float64
	Period float64
}

type WavesService struct {
	api InfoProvider[Waves]
	pub InfoPublisher[Waves]
}

func NewWavesService(api InfoProvider[Waves], pub InfoPublisher[Waves]) *WavesService {
	return &WavesService{
		api: api,
		pub: pub,
	}
}

func (s *WavesService) GetWavesReport(ctx context.Context, city string) (*InfoDataProvider[Waves], error) {
	return s.api.FetchWeatherReport(ctx, city)
}
