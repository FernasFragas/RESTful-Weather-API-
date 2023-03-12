package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"net/url"
)

const key = "6dae91b720b11ea188190cfe41708199"

type WeatherData struct {
	//TODO
}

func GetWeather(ctx *fiber.Ctx) error {
	//ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("q", "Lisbon")
	queryParams.Add("units", "metric")
	queryParams.Add("APPID", key)

	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?%s", queryParams.Encode())

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body) //in order to avoid leaks

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(r))
	err = ctx.Send(r)
	if err != nil {
		return err
	}

	fmt.Sprintf("%s", r)

	return nil
}

// processData
func processData(data []byte) {

}
