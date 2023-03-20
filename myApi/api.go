package myApi

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// GetWeather retrieves the weather data for the city
// Parameters:
//
//	ctx: pointer to the fiber context object
//	city: string representing the name of the city to get the weather for
//
// Returns:
//
//	A string containing the weather information for the given city, or an error if there was a problem
func GetWeather(ctx *fiber.Ctx, city string, weather *WeatherData) error {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("WEATHER_API_KEY")

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("q", city)
	queryParams.Add("units", "metric")
	queryParams.Add("APPID", key)

	// creates a string with the url used to do http request
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?%s", queryParams.Encode())

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	// this is an anonymous function that takes a single parameter of type io.ReadCloser, named Body
	defer func(Body io.ReadCloser) { // schedules the enclosed function to be called at the end of the current function
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
		// closes the anonymous function and passes the resp.Body variable as the parameter for the Body parameter in the function
	}(resp.Body)

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ctx.Send(r)
	if err != nil {
		return err
	}

	marshalError := json.Unmarshal(r, &weather)
	if marshalError != nil {
		return marshalError
	}
	return nil
}
