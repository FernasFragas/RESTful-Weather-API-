package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"restful_api_weather"
)

const port = ":8080"

func main() {

	var key string

	LoadEnvKey(&key, "local.env")

	weatherService := restful_api_weather.NewWeatherService(restful_api_weather.NewWeatherAPI(key), nil /*console to log*/)

	server := restful_api_weather.NewAppServer(weatherService)

	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadEnvKey(key *string, path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	*key = os.Getenv("WEATHER_API_KEY")
}
