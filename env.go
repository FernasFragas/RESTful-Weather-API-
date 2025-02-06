package restful_api_weather

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvKey(key *string, path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	*key = os.Getenv("WEATHER_API_KEY")
}
