package main

import (
	"log"
	"weatherservice"
	"weatherservice/api"
)

const port = ":8080"

func main() {

	weatherServiceSecrets := weatherservice.LoadEnvKey()

	reporters := weatherservice.NewReporters(
		api.NewWeatherAPI(weatherServiceSecrets.OpenWeatherAPIKey),
		api.NewOpenMateoAPI(),
	)

	server := weatherservice.NewAppServer(reporters)

	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
