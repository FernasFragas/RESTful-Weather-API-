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
		city := ctx.FormValue("CityName")
		data, err := myApi.GetWeather(ctx, city)
		if err != nil {
			return err
		}
		return ctx.Render("index", data)
	})

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err.Error())
	}
}
