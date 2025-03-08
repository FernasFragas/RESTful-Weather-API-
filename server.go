package weatherservice

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Reporter interface {
	GenerateReport(ctx context.Context, localization string) (*GeneralInfo, error)
}

type Server struct {
	app *fiber.App

	reporters *Reporters
}

func NewAppServer(reporters *Reporters) *Server {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".go.tpl"),
	})

	server := &Server{
		app:       app,
		reporters: reporters,
	}

	// Serve static files from the "public" directory
	app.Static("/", "./public")

	app.Get("/", server.listGeneralInfo)

	app.Post("/process-form/:CityName", server.listGeneralInfo)

	return server
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(port)
}

func (s *Server) listGeneralInfo(ctx *fiber.Ctx) error {
	city := ctx.FormValue("city_name") // retrieves the name passed in the form
	if city == "" {
		city = "Lisbon"
	}

	generalInfo, err := s.reporters.GenerateReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)

	return ctx.Render("index", map[string]any{"GeneralInfo": generalInfo}) // renders the weather.html file with the data retrieved from the API

}
