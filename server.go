package weatherservice

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Reporter[T any] interface {
	GetReport(ctx context.Context, city string) (*ReportData[T], error)
}

type Server struct {
	app *fiber.App

	weatherReporter *WeatherReporter
	wavesReporter   *WavesReporter
}

func NewAppServer(wr *WeatherReporter, wvr *WavesReporter) *Server {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".go.tpl"),
	})

	server := &Server{
		app:             app,
		weatherReporter: wr,
		wavesReporter:   wvr,
	}

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	app.Post("/process-form/:CityName", server.reportWeather)

	return server
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(port)
}

func (s *Server) reportWeather(ctx *fiber.Ctx) error {
	city := ctx.FormValue("city_name") // retrieves the name passed in the form

	if city == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("city name is empty")
	}

	weatherdt, err := s.weatherReporter.GetReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)

	return ctx.Render("index", map[string]interface{}{"Weather": weatherdt}) // renders the weather.html file with the data retrieved from the API
}

func (s *Server) reportWaves(ctx *fiber.Ctx) error {
	city := ctx.FormValue("city_name") // retrieves the name passed in the form

	if city == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("city name is empty")
	}

	wavesdt, err := s.wavesReporter.GetReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)

	return ctx.Render("index", map[string]interface{}{"Waves": wavesdt}) // renders the waves.html file with the data retrieved from the API
}
