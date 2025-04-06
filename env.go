package weatherservice

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type WeatherServiceKeys struct {
	OpenWeatherAPIKey  string
	StormGlassAPIKey   string
	MateoMaticsAuths   MateoMaticsSecrets
	YoutubeAPIKey      string
	MarkcorpsAPIKey    string
	AmadeusID          string
	AmadeusSecret      string
	GooglePlacesAPIKey string
}

type MateoMaticsSecrets struct {
	Username string
	Password string
}

func LoadEnvKey() (weatherServiceKeys *WeatherServiceKeys) {
	weatherServiceKeys = &WeatherServiceKeys{}

	// Load .env file only in development
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file", err.Error())
		}
	}

	weatherServiceKeys.OpenWeatherAPIKey = os.Getenv("WEATHER_API_KEY")

	//weatherServiceKeys.StormGlassAPIKey = os.Getenv("STORMGLASS_API_KEY")

	//weatherServiceKeys.MateoMaticsAuths.Username = os.Getenv("MATEOMATICS_USERNAME")
	//	weatherServiceKeys.MateoMaticsAuths.Password = os.Getenv("MATEOMATICS_PASSWORD")

	weatherServiceKeys.YoutubeAPIKey = os.Getenv("YOUTUBE_API_KEY")

	//	weatherServiceKeys.MarkcorpsAPIKey = os.Getenv("MARKCORPS_API_KEY")

	//weatherServiceKeys.AmadeusID = os.Getenv("AMADEUS_ID")
	//weatherServiceKeys.AmadeusSecret = os.Getenv("AMADEUS_CLIENT_SECRET")

	weatherServiceKeys.GooglePlacesAPIKey = os.Getenv("GOOGLE_PLACES_API_KEY")

	return weatherServiceKeys
}
