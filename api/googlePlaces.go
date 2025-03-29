package api

import (
	"context"
	"net/http"
	"net/url"
	"weatherservice"
)

const googlePlaceMapTextSearchURL = "https://maps.googleapis.com/maps/api/place/textsearch/json?"

type GooglePlacesAPI struct {
	client *http.Client

	key string
}

func NewGooglePlacesAPI(key string) *GooglePlacesAPI {
	return &GooglePlacesAPI{
		client: http.DefaultClient,
		key:    key,
	}
}

func (api *GooglePlacesAPI) FetchReportData(ctx context.Context, hotelName string) (weatherservice.DataToReport[weatherservice.HotelPhotos], error) {
	hotelSearch := api.searchHotel(ctx, hotelName)

	params := url.Values{}
	params.Add("place_id", placeId)
	params.Add("key", api.key)

	url := googlePlaceMapTextSearchURL + params.Encode()

}

func (api *GooglePlacesAPI) searchHotel(ctx context.Context, hotelName string) (string, error) {
	params := url.Values{}
	params.Add("query", hotelName)
	params.Add("key", api.key)

	url := googlePlaceMapTextSearchURL + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

}
