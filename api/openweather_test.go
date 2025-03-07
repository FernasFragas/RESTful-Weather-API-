package api

import (
	"log"
	"testing"
	"weatherservice"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_LoadEnvKey(t *testing.T) {
	key := weatherservice.LoadEnvKey()

	assert.NotEmpty(t, key)
}

func Test_GetWeather(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	key := weatherservice.LoadEnvKey()

	op := NewWeatherAPI(key.OpenWeatherAPIKey)
	report, err := op.FetchReportData(ctx.Context(), "Lisbon")
	if err != nil {
		log.Fatal("Error ", err.Error())
	}

	assert.NotNil(t, report)
	assert.NotNil(t, report.Data)
	assert.Greater(t, report.Data.Temperature, float64(0))
	assert.Greater(t, report.Data.Humidity, float64(0))
	assert.NotEmpty(t, report.Data.Condition)
}
