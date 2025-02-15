package main

import (
	"log"
	"weatherservice"
)

const port = ":8080"

func main() {

	weatherServiceSecrets := weatherservice.LoadEnvKey()

	weatherService := weatherservice.NewWeatherReporter(weatherservice.NewWeatherAPI(weatherServiceSecrets.OpenWeatherAPIKey), nil)

	server := weatherservice.NewAppServer(weatherService, nil)

	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
