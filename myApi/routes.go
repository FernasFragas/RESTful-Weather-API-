package myApi

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func SetRoutes(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	app.Post("/process-form/:CityName", func(ctx *fiber.Ctx) error {
		city := ctx.FormValue("CityName") // retrieves the name passed in the form
		var weather WeatherData
		err := GetWeather(ctx, city, &weather)
		if err != nil {
			return err
		}
		log.Println(weather)
		return ctx.Render("index", weather) // links the data retrieved from the API and connects it to the html file
	})

}
