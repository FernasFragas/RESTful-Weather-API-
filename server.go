package restful_api_weather

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Server struct {
	app *fiber.App

	service *WeatherService
}

func NewAppServer(s *WeatherService) *Server {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	server := &Server{
		app:     app,
		service: s,
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

	weatherdt, err := s.service.GetWeatherReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	return ctx.Render("weather", weatherdt) // renders the weather.html file with the data retrieved from the API
}
