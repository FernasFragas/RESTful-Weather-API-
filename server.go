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

	app.Get("/videos", server.listVideoStreamInfo)

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

func (s *Server) listGeneralInfo(ctx *fiber.Ctx) error {
	city := ctx.FormValue("city_name") // retrieves the name passed in the form
	if city == "" {
		city = "Lisbon"
	}

	generalInfo, err := s.weatherReporters.GenerateReport(ctx.Context(), city)
	if err != nil {
		return err
	}

	// --- BEGIN Cache Check ---
	cachedJSON, err := GetCityData(city)
	if err != nil && err != sql.ErrNoRows {
		// Handle potential DB errors (other than not found)
		log.Printf("Error checking cache for city %s: %v", city, err)
		// Decide how to handle this - maybe proceed to fetch fresh data, or return an error
		// For now, let's proceed to fetch fresh data, but log the error.
	} else if err == nil {
		// Cache hit!
		log.Printf("Cache hit for city: %s", city)
		var data map[string]any
		if unmarshalErr := json.Unmarshal([]byte(cachedJSON), &data); unmarshalErr != nil {
			log.Printf("Error unmarshaling cached data for city %s: %v", city, unmarshalErr)
			// Data in DB is corrupted? Proceed to fetch fresh data.
		} else {
			data["GeneralInfo"] = generalInfo

			// Successfully got data from cache
			// Need to update session if necessary? Currently, session seems mostly for videos/generalInfo separately.
			// Let's keep it simple and just render the cached data for now.
			if ctx.Get("HX-Request") == "true" {
				return ctx.Render("content_fragment", data)
			}
			return ctx.Render("index", data)
		}
	}
	// --- END Cache Check (If cache miss or error, continue below) ---

	log.Printf("Cache miss for city: %s. Fetching fresh data.", city)

	ctx.Status(fiber.StatusOK)

	sess := ctx.Locals("session").(*session.Session)
	sess.Set("GeneralInfo", generalInfo)

	videos, err := s.videoStreamReporters.GenerateReport(ctx.Context(), fmt.Sprintf("Turistic places in %s, %s", city, generalInfo.Country))
	if err != nil {
		videos = &VideosStream{}
	}

	err = sess.Save()
	if err != nil {
		// Consider logging this error but maybe not returning it to the client
		fmt.Println("Session save error:", err)
	}

	hotels, err := s.hotelsApi.GenerateReport(ctx.Context(), fmt.Sprintf("%f,%f", generalInfo.Lat, generalInfo.Lon))
	if err != nil {
		// Handle hotel error appropriately, maybe return an empty list or log
		fmt.Println("Error fetching hotels:", err)
		hotels = &Hotels{}
	} else {
		fmt.Printf("Fetched %d hotels\n", len(*hotels))
	}

	data := map[string]any{
		"GeneralInfo": generalInfo,
		"Videos":      videos,
		"Hotels":      hotels,
	}

	// Save the fetched data to the database
	go func(cityToSave string, dataToSave map[string]any) {
		if err := SaveCityData(cityToSave, dataToSave); err != nil {
			log.Printf("Error saving data for city %s to DB: %v", cityToSave, err)
		}
	}(city, data) // Pass copies to the goroutine

	// Check if it's an HTMX request
	if ctx.Get("HX-Request") == "true" {
		// Render only the content fragment for HTMX requests
		return ctx.Render("content_fragment", data)
	}

	// Render the full page for regular requests
	return ctx.Render("index", data)
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
