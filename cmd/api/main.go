package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

const port = ":8080"

func main() {
	app := fiber.New()

	app.Get("/", GetWeather)

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err.Error())
	}
}
