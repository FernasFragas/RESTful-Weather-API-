package weatherservice

import (
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
	itenaryReporter      Reporter[Coordinates]
}

type TemplateData struct {
	Query          string
	Videos         VideosStream
	GeneralInfo    GeneralWeatherInfo
	Hotels         Hotels
	ItineraryItems *Coordinates
}

var store = session.New()

func init() {
	// Register the VideosStream type with gob
	gob.Register(VideosStream{})
	gob.Register(GeneralWeatherInfo{})
}

func NewAppServer(weatherReporters Reporter[GeneralWeatherInfo], videoStreamReporters Reporter[VideosStream], hotelsApi Reporter[Hotels], itenaryReporter Reporter[Coordinates]) *Server {
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
		itenaryReporter:      itenaryReporter,
	}

	// Serve static files from the "public" directory
	app.Static("/", "./public")

	app.Get("/", server.listGeneralInfo)

	app.Get("/process-form/", server.listGeneralInfo)

	app.Get("/api/countries", server.getCountries)

	app.Get("/api/cities", server.getCitiesByCountry)

	app.Post("/generate-itinerary", server.generateItinerary)

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
	cityAndCountry := ctx.FormValue("city_name") // retrieves the name passed in the form
	if cityAndCountry == "" {
		cityAndCountry = "Lisbon, Portugal"
	}

	generalInfo, err := s.weatherReporters.GenerateReport(ctx.Context(), cityAndCountry)
	if err != nil {
		log.Printf("Error Retriving New Weather data information", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	city := generalInfo.City

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

	if city == "" {
		return TemplateData{}, fmt.Errorf("city is required")
	}

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

func (s *Server) generateItinerary(ctx *fiber.Ctx) error {
	// Get form data
	city := ctx.FormValue("city_itenary")
	startDate := ctx.FormValue("start_date")
	endDate := ctx.FormValue("end_date")

	// Get selected categories (multiple values)
	// For checkboxes with the same name, we need to get all values
	var categories []string
	formData := ctx.Context().PostArgs()
	formData.VisitAll(func(key, value []byte) {
		if string(key) == "categories" {
			categories = append(categories, string(value))
		}
	})

	// Log the received data for debugging
	log.Printf("Generating itinerary for city: %s", city)
	log.Printf("Start date: %s, End date: %s", startDate, endDate)
	log.Printf("Selected categories: %v", categories)

	// Validate required fields
	if city == "" || startDate == "" || endDate == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "City, start date, and end date are required",
		})
	}

	// Fetch general info to get lon/lat coordinates
	var generalInfo TemplateData
	if data, err := s.checkDatabase(city); err != nil {
		log.Printf("Error checking database for city %s while generating itinerary: %v", city, err)
		log.Printf("Fetching fresh data for city %s", city)

		data, err := s.weatherReporters.GenerateReport(ctx.Context(), city)
		if err != nil {
			log.Printf("Error fetching general info for city %s: %v", city, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch city information",
			})
		}

		generalInfo = TemplateData{
			GeneralInfo: *data,
		}
	} else {
		generalInfo = data
	}

	if len(categories) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one category must be selected",
		})
	}

	// Here you can use your ItenaryReporter to generate the itinerary
	// For now, let's create a simple response and log it
	log.Printf("Itinerary data: city=%s, startDate=%s, endDate=%s, categories=%v",
		city, startDate, endDate, categories)

	categoriesToSearch := fmt.Sprintf("%f,%f", generalInfo.GeneralInfo.Lon, generalInfo.GeneralInfo.Lat) + "," + strings.Join(categories, ",")

	itenary, err := s.itenaryReporter.GenerateReport(ctx.Context(), categoriesToSearch)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate itinerary",
		})
	}

	// Check if it's an HTMX request
	if ctx.Get("HX-Request") == "true" {
		// Return only the itinerary section for HTMX requests
		return ctx.Render("itinerary_card", fiber.Map{
			"Query":          city,
			"GeneralInfo":    generalInfo.GeneralInfo,
			"ItineraryItems": itenary,
		})
	}

	// For non-HTMX requests, redirect back to main page
	return ctx.Redirect("/?city_name=" + city)
}

func (s *Server) getCountries(ctx *fiber.Ctx) error {
	countries, err := GetAllCountries()
	if err != nil {
		log.Printf("Error getting countries: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get countries",
		})
	}

	return ctx.JSON(fiber.Map{
		"countries": countries,
	})
}

func (s *Server) getCitiesByCountry(ctx *fiber.Ctx) error {
	country := ctx.Query("country")
	if country == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Country parameter is required",
		})
	}

	cities, err := GetCitiesByCountry(country)
	if err != nil {
		log.Printf("Error getting cities for country %s: %v", country, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get cities",
		})
	}

	return ctx.JSON(fiber.Map{
		"cities": cities,
	})
}
