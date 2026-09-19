package weatherservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadEnvKey(t *testing.T) {
	// Production mode skips .env, so the test only sees the values set here.
	t.Setenv("ENV", "production")
	t.Setenv("WEATHER_API_KEY", "weather-key")
	t.Setenv("YOUTUBE_NEW", "youtube-key")
	t.Setenv("PLACES_API_NEW", "places-key")
	t.Setenv("GEOAPIFY", "geoapify-key")
	t.Setenv("FOURSQUARE_API_KEY", "foursquare-key")

	keys := LoadEnvKey()

	assert.Equal(t, "weather-key", keys.OpenWeatherAPIKey)
	assert.Equal(t, "youtube-key", keys.YoutubeAPIKey)
	assert.Equal(t, "places-key", keys.GooglePlacesAPIKey)
	assert.Equal(t, "geoapify-key", keys.GeoapifyAPIKey)
	assert.Equal(t, "foursquare-key", keys.FoursquareAPIKey)
}
