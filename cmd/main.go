package main

import (
	"RESTful-Weather-API-/myApi"
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

	// setting app routes
	myApi.SetRoutes(app)

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err.Error())
	}
}
