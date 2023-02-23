package RESTful_Weather_API_

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetWeather(ctx *fiber.Ctx) error {
	//ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	client := &http.Client{}

	queryParams := url.Values{}
	queryParams.Add("datasetid", "GHCND")
	queryParams.Add("locationid", "CITY:PO0001")
	queryParams.Add("datatypeid", "TAVG")
	queryParams.Add("startdate", "2023-02-01")
	queryParams.Add("enddate", "2023-02-23")

	req, err := http.NewRequest("GET", "https://www.ncdc.noaa.gov/cdo-web/data?", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Set("token", "rPcrJuGpkVueOJWXdUfoeRweXVVAbLUT")
	req.URL.RawQuery = queryParams.Encode()
	println(req.URL.Path)
	println(req.URL.RawQuery)
	//datasetid=GHCND&locationid=ZIP:28801&startdate=2010-05-01&enddate=2010-05-01
	//req.Header.Set("datasetid", "GHCND")
	//req.Header.Set("locationid", "ZIP:28801")
	//req.Header.Set("startdate", time.Now().Format("02-Jan-2006"))
	//req.Header.Set("enddate", time.Now().Format("02-Jan-2006"))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close() //in order to avoid leaks

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(r))
	ctx.Send(r)
	return nil
}
