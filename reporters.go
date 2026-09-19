package weatherservice

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

const weatherEmbedURL = "https://embed.waze.com/iframe?zoom=10&lat=%f&lon=%f"

type ReporterProvider[T any] interface {
	FetchReportData(ctx context.Context, _ ...string) (*DataToReport[T], error)
	FetchGeneralInfo(ctx context.Context, _ ...string) (*DataToReport[T], error)
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
		EmbedURL: fmt.Sprintf(weatherEmbedURL, weatherInfo.Data.Lat, weatherInfo.Data.Lon),
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

type HotelsApi struct {
	hotelsApi ReporterProvider[Hotels]
	photosAPI ReporterProvider[HotelsPhotos]
}

func NewHotelsApi(hotelsApi ReporterProvider[Hotels], photosAPI ReporterProvider[HotelsPhotos]) *HotelsApi {
	return &HotelsApi{
		hotelsApi: hotelsApi,
		photosAPI: photosAPI,
	}
}

type Hotels []Hotel

type Hotel struct {
	HotelName    string
	HotelURL     string
	HotelPrice   string
	HotelRating  float64
	HotelAddress string
	HotelMapURL  string
	ContactPhone string
	PriceRange   string
	HotelPhotos  []string
	HotelReviews []HotelReview
}

type HotelReview struct {
	AuthorName string
	Text       string
	Rating     float64
}

func (s *HotelsApi) GenerateReport(ctx context.Context, city string) (*Hotels, error) {
	hotels, err := s.hotelsApi.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	data := make(Hotels, len(hotels.Data))

	for i, hotel := range hotels.Data {
		hotelsPhotos := make(chan []string, len(hotel.HotelPhotos))

		wg := sync.WaitGroup{}
		for _, photo := range hotel.HotelPhotos {
			wg.Add(1)
			go func(pic string) {
				defer wg.Done()
				if err := s.retrieveHotelPhotos(ctx, pic, hotelsPhotos); err != nil {
					log.Printf("Error retrieving photos for %s: %v", pic, err)
				}
			}(photo)
		}

		go func() {
			wg.Wait()
			close(hotelsPhotos)
		}()

		for photos := range hotelsPhotos {
			data[i].HotelPhotos = append(data[i].HotelPhotos, photos...)
		}

		data[i].HotelName = hotel.HotelName
		data[i].HotelURL = hotel.HotelURL
		data[i].HotelPrice = hotel.HotelPrice
		data[i].HotelRating = hotel.HotelRating
		data[i].HotelAddress = hotel.HotelAddress
		data[i].HotelMapURL = hotel.HotelMapURL
		data[i].ContactPhone = hotel.ContactPhone
		data[i].PriceRange = hotel.PriceRange
		data[i].HotelReviews = hotel.HotelReviews
	}

	return &data, nil
}

func (s *HotelsApi) retrieveHotelPhotos(ctx context.Context, hotelName string, hotelsPhotos chan []string) error {
	photo, err := s.photosAPI.FetchReportData(ctx, hotelName)
	if err != nil {
		return err
	}

	hotelsPhotos <- photo.Data[0].HotelPhotos

	return nil
}

type HotelsPhotos []HotelPhoto

type HotelPhoto struct {
	HotelName   string
	HotelPhotos []string
	HotelURL    string
}

type FlightsInfo []FlightInfo

type FlightInfo struct {
	DepartureDate  string
	ArrivalDate    string
	Airline        string
	FlightNumber   string
	FlightDuration string
	Price          float64
	Origin         string
	Destination    string
	FlightMapURL   string
}

// Geocoding API Response Structs
type GeocodingResponse []GeocodingResult

type GeocodingResult struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names,omitempty"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state,omitempty"`
}

type CoordinatesReporter struct {
	api ReporterProvider[Coordinates]
}

func NewCoordinatesReporter(api ReporterProvider[Coordinates]) *CoordinatesReporter {
	return &CoordinatesReporter{
		api: api,
	}
}

type Itenary struct {
	Categories        []string
	PlacesCoordinates []Coordinates
	StartDate         string
	EndDate           string
	Location          string
	Description       string
	ImageURL          string
	MapURL            string
}

type Coordinates struct {
	CoordinatesWithName []CoordinatesWithName
}

type CoordinatesWithName struct {
	Name             string
	Coordinates      []interface{}
	CoordinatesFloat []float64
}

func (s *CoordinatesReporter) GenerateReport(ctx context.Context, city string) (*Coordinates, error) {
	coordinates, err := s.api.FetchReportData(ctx, city)
	if err != nil {
		return nil, err
	}

	coordinatesWithName := make([]CoordinatesWithName, len(coordinates.Data.CoordinatesWithName))

	for i, coord := range coordinates.Data.CoordinatesWithName {
		coordinatesFloat, err := parseCoordinates(coord.Coordinates)
		if err != nil {
			return nil, err
		}

		coordinatesWithName[i].CoordinatesFloat = coordinatesFloat
		coordinatesWithName[i].Name = coord.Name
		coordinatesWithName[i].Coordinates = coord.Coordinates
	}

	return &Coordinates{
		CoordinatesWithName: coordinatesWithName,
	}, nil
}

func parseCoordinates(coordinates []interface{}) ([]float64, error) {
	coordinatesFloat := make([]float64, len(coordinates))

	for i, coord := range coordinates {
		coordFloat, ok := coord.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid coordinate type: %T", coord)
		}
		coordinatesFloat[i] = coordFloat
	}

	return coordinatesFloat, nil
}
