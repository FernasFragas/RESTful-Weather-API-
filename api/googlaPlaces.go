package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"weatherservice"
)

const googleNearByCoordinationSearchURL = "https://places.googleapis.com/v1/places:searchNearby"

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

func (api *GooglePlacesAPI) FetchReportData(ctx context.Context, CityName ...string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {
	coord := strings.Split(CityName[0], ",")

	latitude, err := strconv.ParseFloat(coord[0], 64)
	if err != nil {
		return nil, err
	}

	longitude, err := strconv.ParseFloat(coord[1], 64)
	if err != nil {
		return nil, err
	}

	hotelSearch, err := api.searchHotels(ctx, latitude, longitude)
	if err != nil {
		return nil, err
	}

	hotels := make(weatherservice.Hotels, len(hotelSearch.Places))

	for i, place := range hotelSearch.Places {
		photos := make([]string, len(place.Photos))
		for j, photo := range place.Photos {
			photos[j] = photo.Name
		}

		reviews := make([]weatherservice.HotelReview, len(place.Reviews))
		for j, review := range place.Reviews {
			reviews[j] = weatherservice.HotelReview{
				AuthorName: review.AuthorAttribution.DisplayName,
				Text:       review.Text.Text,
				Rating:     review.Rating,
			}
		}

		hotels[i] = weatherservice.Hotel{
			HotelName:    place.DisplayName.Text,
			HotelURL:     place.WebsiteUri,
			HotelRating:  place.Rating,
			HotelAddress: place.FormattedAddress,
			HotelMapURL:  place.GoogleMapsUri,
			ContactPhone: place.InternationalPhoneNumber,
			HotelPhotos:  photos,
			HotelReviews: reviews,
		}
	}

	return &weatherservice.DataToReport[weatherservice.Hotels]{
		Data: hotels,
	}, nil
}

type nearbySearchRequest struct {
	MaxResultCount      int                 `json:"maxResultCount"`
	IncludedTypes       []string            `json:"includedTypes"`
	ExcludedTypes       []string            `json:"excludedTypes"`
	RankPreference      string              `json:"rankPreference"`
	LocationRestriction locationRestriction `json:"locationRestriction"`
}

type locationRestriction struct {
	Circle circle `json:"circle"`
}

type circle struct {
	Center center  `json:"center"`
	Radius float64 `json:"radius"`
}

type center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (api *GooglePlacesAPI) searchHotels(ctx context.Context, latitude, longitude float64) (*PlacesResponse, error) {
	searchReq := nearbySearchRequest{
		IncludedTypes:  []string{"hotel"}, // Matches the curl command exactly
		MaxResultCount: 8,
		RankPreference: "POPULARITY",
		LocationRestriction: locationRestriction{
			Circle: circle{
				Center: center{
					Latitude:  latitude,
					Longitude: longitude,
				},
				Radius: 5000.0, // Matches the radius in the curl command
			},
		},
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", googleNearByCoordinationSearchURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", api.key)
	req.Header.Set("X-Goog-FieldMask", "places")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Check if response is not successful
	if resp.StatusCode != http.StatusOK {
		// Read the error response body
		var errorBody []byte
		errorBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP %d: failed to read error body", resp.StatusCode)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(errorBody))
	}

	var response *PlacesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type PlacesResponse struct {
	Places []PlaceHotel `json:"places"`
}

type PlaceHotel struct {
	Name                     string         `json:"name"`
	ID                       string         `json:"id"`
	Types                    []string       `json:"types"`
	DisplayName              DisplayName    `json:"displayName"`
	Location                 LatLng         `json:"location"`
	Rating                   float64        `json:"rating,omitempty"`
	Photos                   []Photo        `json:"photos,omitempty"`
	FormattedAddress         string         `json:"formattedAddress,omitempty"`
	WebsiteUri               string         `json:"websiteUri,omitempty"`
	GoogleMapsUri            string         `json:"googleMapsUri,omitempty"`
	InternationalPhoneNumber string         `json:"internationalPhoneNumber,omitempty"`
	Reviews                  []ReviewPlaces `json:"reviews,omitempty"`
	UserRatingCount          int            `json:"userRatingCount,omitempty"`
	IconMaskBaseUri          string         `json:"iconMaskBaseUri,omitempty"`
	IconBackgroundColor      string         `json:"iconBackgroundColor,omitempty"`
}

type DisplayName struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type ReviewPlaces struct {
	Text              LocalizedText     `json:"text"`
	Rating            float64           `json:"rating"`
	AuthorAttribution AuthorAttribution `json:"authorAttribution"`
}

func (api *GooglePlacesAPI) FetchGeneralInfo(ctx context.Context, city ...string) (*weatherservice.DataToReport[weatherservice.Hotels], error) {
	return nil, nil
}
