package weatherservice

import "context"

type ReportData[T any] struct {
	Data T
}

type ReporterProvider[T any] interface {
	FetchReportData(ctx context.Context, city string) (*ReportData[T], error)
}

type ReporterPublisher[T any] interface {
	PublishReportData(ctx context.Context, weather *ReportData[T]) error
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

type WeatherReporter struct {
	api ReporterProvider[Weather]
	pub ReporterPublisher[Weather]
}

func NewWeatherReporter(api ReporterProvider[Weather], pub ReporterPublisher[Weather]) *WeatherReporter {
	return &WeatherReporter{
		api: api,
		pub: pub,
	}
}

func (s *WeatherReporter) GetReport(ctx context.Context, city string) (*ReportData[Weather], error) {
	return s.api.FetchReportData(ctx, city)
}

type Waves struct {
	Height float64
	Period float64
}

type WavesReporter struct {
	api ReporterProvider[Waves]
	pub ReporterPublisher[Waves]
}

func NewWavesReporter(api ReporterProvider[Waves], pub ReporterPublisher[Waves]) *WavesReporter {
	return &WavesReporter{
		api: api,
		pub: pub,
	}
}

func (s *WavesReporter) GetReport(ctx context.Context, city string) (*ReportData[Waves], error) {
	return s.api.FetchReportData(ctx, city)
}
