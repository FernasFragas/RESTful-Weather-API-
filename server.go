package weatherservice

import (
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

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

	app.Get("/process-form/", server.listGeneralInfo)

	return server
}

// Initialize Database before starting server
func (s *Server) InitializeDatabase(dbPath string) error {
	return InitDB(dbPath)
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(port)
}

type TemplateData struct {
	Query       string
	Videos      VideosStream
	GeneralInfo GeneralWeatherInfo
	Hotels      Hotels
}

func (s *Server) listGeneralInfo(ctx *fiber.Ctx) error {
	city := ctx.FormValue("city_name") // retrieves the name passed in the form
	if city == "" {
		city = "Lisbon"
	}

	generalInfo, err := s.weatherReporters.GenerateReport(ctx.Context(), city)
	if err != nil {
		log.Printf("Error Retriving New Weather data information", err)
	}

	// --- BEGIN Cache Check ---
	var data TemplateData

	if data, err = s.checkDatabase(city); err != nil {
		log.Printf("Error Retriving search data information with error %s", err)

		log.Printf("Cache miss for city: %s. Fetching fresh data.", city)

		data, err = s.retireveFreshInformation(ctx, generalInfo, city)
		if err != nil {
			log.Printf("Error Retriving fresh data information with error %s", err)
		}
	}

	// Check if it's an HTMX request
	if ctx.Get("HX-Request") == "true" {
		// Render only the content fragment for HTMX requests
		return ctx.Render("content_fragment", fiber.Map{
			"Query":       city,
			"GeneralInfo": data.GeneralInfo,
			"Videos":      data.Videos,
			"Hotels":      data.Hotels,
		})
	}

	// Render the full page for regular requests
	return ctx.Render("index", fiber.Map{
		"Query":       city,
		"GeneralInfo": data.GeneralInfo,
		"Videos":      data.Videos,
		"Hotels":      data.Hotels,
	})
}

func (s *Server) checkDatabase(city string) (TemplateData, error) {
	var data TemplateData

	cachedJSON, err := GetCityData(city)
	if err != nil && err != sql.ErrNoRows {
		// Handle potential DB errors (other than not found)
		log.Printf("Error checking cache for city %s: %v", city, err)
		// Decide how to handle this - maybe proceed to fetch fresh data, or return an error
		// For now, let's proceed to fetch fresh data, but log the error.
	} else if err == nil {
		// Cache hit!
		log.Printf("Cache hit for city: %s", city)

		if unmarshalErr := json.Unmarshal([]byte(cachedJSON), &data); unmarshalErr != nil {
			log.Printf("Error unmarshaling cached data for city %s: %v", city, unmarshalErr)
			// Data in DB is corrupted? Proceed to fetch fresh data.
		}
	} else {
		return TemplateData{}, err
	}

	return data, nil
}

func (s *Server) retireveFreshInformation(ctx *fiber.Ctx, generalInfo *GeneralWeatherInfo, city string) (TemplateData, error) {
	ctx.Status(fiber.StatusOK)

	videos, err := s.videoStreamReporters.GenerateReport(ctx.Context(), fmt.Sprintf("Turistic places in %s, %s", city, generalInfo.Country))
	if err != nil {
		videos = &VideosStream{}
	}

	hotels, err := s.hotelsApi.GenerateReport(ctx.Context(), fmt.Sprintf("%f,%f", generalInfo.Lat, generalInfo.Lon))
	if err != nil {
		// Handle hotel error appropriately, maybe return an empty list or log
		fmt.Println("Error fetching hotels:", err)

		hotels = &Hotels{}
	} else {
		fmt.Printf("Fetched %d hotels\n", len(*hotels))
	}

	data := TemplateData{
		GeneralInfo: *generalInfo,
		Videos:      *videos,
		Hotels:      *hotels,
	}

	// Save the fetched data to the database
	go func(cityToSave string, dataToSave TemplateData) {
		dtToSave := map[string]any{
			"GeneralInfo": dataToSave.GeneralInfo,
			"Videos":      dataToSave.Videos,
			"Hotels":      dataToSave.Hotels,
		}

		if err := SaveCityData(cityToSave, dtToSave); err != nil {
			log.Printf("Error saving data for city %s to DB: %v", cityToSave, err)
		}
	}(city, data)

	if err != nil {
		return data, err
	}

	return data, nil
}
