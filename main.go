package main

import (
	"RESTful-Weather-API-/cmd/myApi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"log"
)

const port = ":8080"

func main() {

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	app.Post("/process-form", func(ctx *fiber.Ctx) error {
		city := ctx.FormValue("CityName") // retrieves the name passed in the form
		data, err := myApi.GetWeather(ctx, city)
		if err != nil {
			return err
		}
		return ctx.Render("index", data) // links the data retrieved from the API and connects it to the html file
	})

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err.Error())
	}
}
