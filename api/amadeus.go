package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"weatherservice"
)

const amadeusBaseURL = "https://test.api.amadeus.com/v1/"
const amadeusAuthURL = "security/oauth2/token"
const amadeusHotelsReportURL = "https://test.api.amadeus.com/v1/reference-data/locations/hotels/by-geocode?"

type AmadeusAPI struct {
	client *http.Client

	AmadeusID     string
	AmadeusSecret string

	token string
}

func NewAmadeusAPI(amadeusID string, amadeusSecret string) *AmadeusAPI {
	api := &AmadeusAPI{
		client:        http.DefaultClient,
		AmadeusID:     amadeusID,
		AmadeusSecret: amadeusSecret,
	}

	token, err := api.fetchAuthToken(context.Background())
	if err != nil {
		return nil
	}

	api.token = token

	return api
}

func (api *AmadeusAPI) FetchReportData(ctx context.Context, city string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {
	if api.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	coord := strings.Split(city, ",")

	lat, err := strconv.ParseFloat(coord[0], 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(coord[1], 64)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	params.Add("radius", "100")
	params.Add("hotelSource", "ALL")
	params.Add("radiusUnit", "KM")

	url := amadeusHotelsReportURL + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+api.token)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var response Response

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	hotels := make(weatherservice.Hotels, len(response.Data))
	for i, hotel := range response.Data {
		hotels[i] = weatherservice.Hotel{
			HotelName: hotel.Name,
			HotelURL:  hotel.HotelID,
		}
	}

	return &weatherservice.DataToReport[weatherservice.Hotels]{Data: hotels}, nil
}

func (api *AmadeusAPI) FetchGeneralInfo(ctx context.Context, _ string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {
	return nil, nil
}

func (api *AmadeusAPI) fetchAuthToken(ctx context.Context) (string, error) {
	if api.client == nil {
		return "", fmt.Errorf("client not initialized")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", api.AmadeusID)
	data.Set("client_secret", api.AmadeusSecret)

	req, err := http.NewRequest("POST", amadeusBaseURL+amadeusAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

type GeoCode struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	CountryCode string `json:"countryCode"`
}

type Distance struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Hotel struct {
	ChainCode  string   `json:"chainCode"`
	IataCode   string   `json:"iataCode"`
	DupeID     int      `json:"dupeId"`
	Name       string   `json:"name"`
	HotelID    string   `json:"hotelId"`
	GeoCode    GeoCode  `json:"geoCode"`
	Address    Address  `json:"address"`
	Distance   Distance `json:"distance"`
	LastUpdate string   `json:"lastUpdate"`
}

type Response struct {
	Data []Hotel `json:"data"`
}
