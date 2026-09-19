package main

import (
	"log"
	"os"
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

	// DB_PATH should point at the mounted volume in production so data survives deploys.
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "weatherservice.db"
	}

	server := weatherservice.NewAppServer(reporters, videoStreamReporters, hotelsApi, coordinatesReporter)
	if err := server.InitializeDatabase(dbPath); err != nil {
		log.Fatal(err)
	}

	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
