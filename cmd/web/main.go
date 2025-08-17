package main

import (
	"log"
	"weatherservice"
	"weatherservice/api"
)

const port = ":8080"

func main() {

	weatherServiceSecrets := weatherservice.LoadEnvKey()
	coordinatesReporter := weatherservice.NewCoordinatesReporter(api.NewFoursquareAPI(weatherServiceSecrets.FoursquareAPIKey))

	reporters := weatherservice.NewWeatherReporters(
		api.NewWeatherAPI(weatherServiceSecrets.OpenWeatherAPIKey),
		api.NewOpenMateoAPI(),
	)

	hotelsApi := weatherservice.NewHotelsApi(
		api.NewGooglePlacesAPI(weatherServiceSecrets.GooglePlacesAPIKey),
		api.NewGooglePhotosAPI(weatherServiceSecrets.GooglePlacesAPIKey),
	)

	videoStreamReporters := weatherservice.NewVideoStreamReporters(
		api.NewYoutubeAPI(weatherServiceSecrets.YoutubeAPIKey),
	)

	server := weatherservice.NewAppServer(reporters, videoStreamReporters, hotelsApi, coordinatesReporter)
	if err := server.InitializeDatabase("weatherservice.db"); err != nil {
		log.Fatal(err)
	}

	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
