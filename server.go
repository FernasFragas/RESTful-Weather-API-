package weatherservice

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

type Reporter[T any] interface {
	GenerateReport(ctx context.Context, localization string) (*T, error)
}

type Server struct {
	app *fiber.App

	weatherReporters     Reporter[GeneralWeatherInfo]
	videoStreamReporters Reporter[VideosStream]
	hotelsApi            Reporter[Hotels]
}

var store = session.New()

func init() {
	// Register the VideosStream type with gob
	gob.Register(VideosStream{})
	gob.Register(GeneralWeatherInfo{})
}

func NewAppServer(weatherReporters Reporter[GeneralWeatherInfo], videoStreamReporters Reporter[VideosStream], hotelsApi Reporter[Hotels]) *Server {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".go.tpl"),
	})

	app.Use(fgprof.New())

	// Use the session middleware
	app.Use(func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}
		ctx.Locals("session", sess)
		return ctx.Next()
	})

	server := &Server{
		app:                  app,
		weatherReporters:     weatherReporters,
		videoStreamReporters: videoStreamReporters,
		hotelsApi:            hotelsApi,
	}

	// Serve static files from the "public" directory
	app.Static("/", "./public")

	app.Get("/", server.listGeneralInfo)

	app.Get("/videos", server.listVideoStreamInfo)

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

	generalInfo, err := s.weatherReporters.GenerateReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)

	sess := ctx.Locals("session").(*session.Session)
	sess.Set("GeneralInfo", generalInfo)

	videos, err := s.videoStreamReporters.GenerateReport(ctx.Context(), fmt.Sprintf("Turistic places in %s, %s", city, generalInfo.Country))
	if err != nil {
		videos = &VideosStream{}
	}

	err = sess.Save()
	if err != nil {
		return err
	}

	hotels, err := s.hotelsApi.GenerateReport(ctx.Context(), fmt.Sprintf("%f,%f", generalInfo.Lat, generalInfo.Lon))
	if err != nil {
		return err
	}

	return ctx.Render("index", map[string]any{"GeneralInfo": generalInfo, "Videos": videos, "Hotels": hotels})
}

func (s *Server) listVideoStreamInfo(ctx *fiber.Ctx) error {
	query := ctx.FormValue("query")
	if query == "" {
		query = "windsurf"
	}

	videos, err := s.videoStreamReporters.GenerateReport(ctx.Context(), query)
	if err != nil {
		return err
	}

	sess := ctx.Locals("session").(*session.Session)
	sess.Set("Videos", videos)

	generalInfo, ok := sess.Get("GeneralInfo").(GeneralWeatherInfo)
	if !ok {
		generalInfo = GeneralWeatherInfo{
			City:    "Lisbon",
			Country: "pt",
			Lon:     -9.1393,
			Lat:     38.7223,
		}
	}

	err = sess.Save()
	if err != nil {
		return err
	}

	return ctx.Render("index", fiber.Map{
		"Query":       query,
		"Videos":      videos,
		"GeneralInfo": generalInfo,
	})
}
