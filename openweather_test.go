package restful_api_weather_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"log"
	"restful_api_weather"
	"testing"
)

func Test_LoadEnvKey(t *testing.T) {
	var key string

	LoadEnvKey(&key, "./../local.env")

	assert.NotEmpty(t, key)
}

func Test_GetWeather(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	var key string

	LoadEnvKey(&key, "./../local.env")

	var data restful_api_weather.WeatherData
	err := restful_api_weather.GetWeather(ctx, "Lisbon", &data, key)
	if err != nil {
		log.Fatal("Error ", err.Error())
	}

	assert.Equal(t, "Lisbon", data.Name)
	assert.Equal(t, "PT", data.Sys.Country)
}
