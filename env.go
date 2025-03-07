package weatherservice

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type WeatherServiceKeys struct {
	OpenWeatherAPIKey string
	StormGlassAPIKey  string
	MateoMaticsAuths  MateoMaticsSecrets
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

	return weatherServiceKeys
}
