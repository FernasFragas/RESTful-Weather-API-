package main

import (
	"log"
	"weatherservice"
	"weatherservice/api"
)

const port = ":8080"

func main() {

	weatherServiceSecrets := weatherservice.LoadEnvKey()

	reporters := weatherservice.NewWeatherReporters(
		api.NewWeatherAPI(weatherServiceSecrets.OpenWeatherAPIKey),
		api.NewOpenMateoAPI(),
	)

	hotelsApi := weatherservice.NewHotelsApi(
		api.NewMarkcorpsAPI(weatherServiceSecrets.MarkcorpsAPIKey),
	)

	videoStreamReporters := weatherservice.NewVideoStreamReporters(
		api.NewYoutubeAPI(weatherServiceSecrets.YoutubeAPIKey),
	)

	server := weatherservice.NewAppServer(reporters, videoStreamReporters, hotelsApi)
	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
