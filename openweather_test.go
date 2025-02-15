package weatherservice

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"log"
	"testing"
)

func Test_LoadEnvKey(t *testing.T) {
	key := LoadEnvKey()

	assert.NotEmpty(t, key)
}

func Test_GetWeather(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	key := LoadEnvKey()

	op := NewWeatherAPI(key.OpenWeatherAPIKey)
	data, err := op.FetchWeatherReport(ctx.Context(), "Lisbon")
	if err != nil {
		log.Fatal("Error ", err.Error())
	}

	assert.Equal(t, "Lisbon", data.City)
	assert.Equal(t, "pt", data.Country)
}
