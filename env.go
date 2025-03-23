package weatherservice

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type WeatherServiceKeys struct {
	OpenWeatherAPIKey string
	StormGlassAPIKey  string
	MateoMaticsAuths  MateoMaticsSecrets
	YoutubeAPIKey     string
	MarkcorpsAPIKey   string
}

type MateoMaticsSecrets struct {
	Username string
	Password string
}

func LoadEnvKey() (weatherServiceKeys *WeatherServiceKeys) {
	weatherServiceKeys = &WeatherServiceKeys{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	weatherServiceKeys.OpenWeatherAPIKey = os.Getenv("WEATHER_API_KEY")

	weatherServiceKeys.StormGlassAPIKey = os.Getenv("STORMGLASS_API_KEY")

	weatherServiceKeys.MateoMaticsAuths.Username = os.Getenv("MATEOMATICS_USERNAME")
	weatherServiceKeys.MateoMaticsAuths.Password = os.Getenv("MATEOMATICS_PASSWORD")

	weatherServiceKeys.YoutubeAPIKey = os.Getenv("YOUTUBE_API_KEY")

	weatherServiceKeys.MarkcorpsAPIKey = os.Getenv("MARKCORPS_API_KEY")

	return weatherServiceKeys
}
