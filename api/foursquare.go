package api

import (
	"context"
	"net/http"
	"weatherservice"
)

const foursquareSearchURL = "https://api.foursquare.com/v3/places/search"

// FoursquareAPI is a client for the Foursquare API.
type FoursquareAPI struct {
	client *http.Client

	token string
}

// NewFoursquareAPI creates a new FoursquareAPI client.
func NewFoursquareAPI(token string) *FoursquareAPI {
	return &FoursquareAPI{
		client: http.DefaultClient,
		token:  token,
	}
}

func (api *FoursquareAPI) FetchReportData(ctx context.Context, city ...string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {

	// need to check if here I receive the city name or the latitude and longitude
	// maybe should change where I gonna get the meteorological data??

	// Create a new HTTP request
	req, err := http.NewRequest("GET", foursquareSearchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", api.token)

	return nil, nil
}

func (api *FoursquareAPI) FetchGeneralInfo(ctx context.Context, city ...string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {
	return nil, nil
}
