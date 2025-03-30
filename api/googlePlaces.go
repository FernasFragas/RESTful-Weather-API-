package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"weatherservice"
)

const googlePlaceMapTextSearchURL = "https://places.googleapis.com/v1/places:searchText"
const googlePhotosURL = "https://places.googleapis.com/v1/%s/media?maxHeightPx=400&maxWidthPx=400&key=%s"

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

func (api *GooglePlacesAPI) FetchReportData(ctx context.Context, hotelName string) (*weatherservice.DataToReport[weatherservice.HotelsPhotos], error) {
	hotelSearch, err := api.searchHotel(ctx, hotelName)
	if err != nil {
		return nil, err
	}

	if len(hotelSearch.Places) == 0 {
		return nil, nil
	}

	hotelPhotos := make(weatherservice.HotelsPhotos, len(hotelSearch.Places))

	for i, place := range hotelSearch.Places {
		if !strings.Contains(strings.ToLower(place.DisplayName.Text), strings.ToLower(hotelName)) {
			continue
		}

		hotelPhotos[i].HotelName = hotelName

		for _, photo := range place.Photos {
			photoId := photo.Name

			params := url.Values{}
			params.Add("key", api.key)

			photoUrl := fmt.Sprintf(googlePhotosURL, photoId, api.key)

			req, err := http.NewRequest("GET", photoUrl, nil)
			if err != nil {
				return nil, err
			}

			resp, err := api.client.Do(req)
			if err != nil {
				return nil, err
			}

			defer func() {
				_ = resp.Body.Close()
			}()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			img := base64.StdEncoding.EncodeToString(body)

			hotelPhotos[i].HotelPhotos = append(hotelPhotos[i].HotelPhotos, img)
		}
	}

	return &weatherservice.DataToReport[weatherservice.HotelsPhotos]{
		Data: hotelPhotos,
	}, nil

	/*
			// Check if the response status is not 200 OK
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("bad request: %s", string(body))
			}

		var photoResponse photoResponse
		err = json.NewDecoder(resp.Body).Decode(&photoResponse)
		if err != nil {
			return nil, err
		}

		return nil, nil*/
}

func (api *GooglePlacesAPI) FetchGeneralInfo(ctx context.Context, _ string) (*weatherservice.DataToReport[weatherservice.HotelsPhotos], error) {
	return nil, nil
}

type SearchRequest struct {
	TextQuery string `json:"textQuery"`
}

func (api *GooglePlacesAPI) searchHotel(ctx context.Context, hotelName string) (*response, error) {
	searchReq := SearchRequest{
		TextQuery: hotelName,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", googlePlaceMapTextSearchURL, bytes.NewBuffer(jsonData))
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

	var response *response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type RoutingSummary struct {
	// Define fields for the RoutingSummary object as needed
}

type ContextualContent struct {
	// Define fields for the ContextualContent object as needed
}

type response struct {
	Places             []Place             `json:"places"`
	RoutingSummaries   []RoutingSummary    `json:"routingSummaries"`
	ContextualContents []ContextualContent `json:"contextualContents"`
	NextPageToken      string              `json:"nextPageToken"`
	SearchUri          string              `json:"searchUri"`
}

type LocalizedText struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type PostalAddress struct {
	// Add fields for postal address as needed
}

type AddressComponent struct {
	// Add fields for address component as needed
}

type PlusCode struct {
	// Add fields for plus code as needed
}

type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Viewport struct {
	// Add fields for viewport as needed
}

type Review struct {
	Name                           string            `json:"name"`
	RelativePublishTimeDescription string            `json:"relativePublishTimeDescription"`
	Text                           LocalizedText     `json:"text"`
	OriginalText                   LocalizedText     `json:"originalText"`
	Rating                         float64           `json:"rating"`
	AuthorAttribution              AuthorAttribution `json:"authorAttribution"`
	PublishTime                    string            `json:"publishTime"`
	FlagContentUri                 string            `json:"flagContentUri"`
	GoogleMapsUri                  string            `json:"googleMapsUri"`
}

type OpeningHours struct {
	// Add fields for opening hours as needed
}

type TimeZone struct {
	// Add fields for time zone as needed
}

type Photo struct {
	Name               string              `json:"name"`
	WidthPx            int                 `json:"widthPx"`
	HeightPx           int                 `json:"heightPx"`
	AuthorAttributions []AuthorAttribution `json:"authorAttributions"`
	FlagContentUri     string              `json:"flagContentUri"`
	GoogleMapsUri      string              `json:"googleMapsUri"`
}

type AuthorAttribution struct {
	DisplayName string `json:"displayName"`
	Uri         string `json:"uri"`
	PhotoUri    string `json:"photoUri"`
}

type Attribution struct {
	// Add fields for attribution as needed
}

type PaymentOptions struct {
	// Add fields for payment options as needed
}

type ParkingOptions struct {
	// Add fields for parking options as needed
}

type SubDestination struct {
	// Add fields for sub destination as needed
}

type FuelOptions struct {
	// Add fields for fuel options as needed
}

type EVChargeOptions struct {
	// Add fields for EV charge options as needed
}

type GenerativeSummary struct {
	// Add fields for generative summary as needed
}

type AreaSummary struct {
	// Add fields for area summary as needed
}

type ContainingPlace struct {
	// Add fields for containing place as needed
}

type AddressDescriptor struct {
	// Add fields for address descriptor as needed
}

type GoogleMapsLinks struct {
	// Add fields for Google Maps links as needed
}

type PriceRange struct {
	// Add fields for price range as needed
}

type AccessibilityOptions struct {
	// Add fields for accessibility options as needed
}

type Place struct {
	Name                         string               `json:"name"`
	ID                           string               `json:"id"`
	DisplayName                  LocalizedText        `json:"displayName"`
	Types                        []string             `json:"types"`
	PrimaryType                  string               `json:"primaryType"`
	PrimaryTypeDisplayName       LocalizedText        `json:"primaryTypeDisplayName"`
	NationalPhoneNumber          string               `json:"nationalPhoneNumber"`
	InternationalPhoneNumber     string               `json:"internationalPhoneNumber"`
	FormattedAddress             string               `json:"formattedAddress"`
	ShortFormattedAddress        string               `json:"shortFormattedAddress"`
	PostalAddress                PostalAddress        `json:"postalAddress"`
	AddressComponents            []AddressComponent   `json:"addressComponents"`
	PlusCode                     PlusCode             `json:"plusCode"`
	Location                     LatLng               `json:"location"`
	Viewport                     Viewport             `json:"viewport"`
	Rating                       float64              `json:"rating"`
	GoogleMapsUri                string               `json:"googleMapsUri"`
	WebsiteUri                   string               `json:"websiteUri"`
	Reviews                      []Review             `json:"reviews"`
	RegularOpeningHours          OpeningHours         `json:"regularOpeningHours"`
	TimeZone                     TimeZone             `json:"timeZone"`
	Photos                       []Photo              `json:"photos"`
	AdrFormatAddress             string               `json:"adrFormatAddress"`
	BusinessStatus               string               `json:"businessStatus"`
	PriceLevel                   string               `json:"priceLevel"`
	Attributions                 []Attribution        `json:"attributions"`
	IconMaskBaseUri              string               `json:"iconMaskBaseUri"`
	IconBackgroundColor          string               `json:"iconBackgroundColor"`
	CurrentOpeningHours          OpeningHours         `json:"currentOpeningHours"`
	CurrentSecondaryOpeningHours []OpeningHours       `json:"currentSecondaryOpeningHours"`
	RegularSecondaryOpeningHours []OpeningHours       `json:"regularSecondaryOpeningHours"`
	EditorialSummary             LocalizedText        `json:"editorialSummary"`
	PaymentOptions               PaymentOptions       `json:"paymentOptions"`
	ParkingOptions               ParkingOptions       `json:"parkingOptions"`
	SubDestinations              []SubDestination     `json:"subDestinations"`
	FuelOptions                  FuelOptions          `json:"fuelOptions"`
	EVChargeOptions              EVChargeOptions      `json:"evChargeOptions"`
	GenerativeSummary            GenerativeSummary    `json:"generativeSummary"`
	AreaSummary                  AreaSummary          `json:"areaSummary"`
	ContainingPlaces             []ContainingPlace    `json:"containingPlaces"`
	AddressDescriptor            AddressDescriptor    `json:"addressDescriptor"`
	GoogleMapsLinks              GoogleMapsLinks      `json:"googleMapsLinks"`
	PriceRange                   PriceRange           `json:"priceRange"`
	UtcOffsetMinutes             int                  `json:"utcOffsetMinutes"`
	UserRatingCount              int                  `json:"userRatingCount"`
	Takeout                      bool                 `json:"takeout"`
	Delivery                     bool                 `json:"delivery"`
	DineIn                       bool                 `json:"dineIn"`
	CurbsidePickup               bool                 `json:"curbsidePickup"`
	Reservable                   bool                 `json:"reservable"`
	ServesBreakfast              bool                 `json:"servesBreakfast"`
	ServesLunch                  bool                 `json:"servesLunch"`
	ServesDinner                 bool                 `json:"servesDinner"`
	ServesBeer                   bool                 `json:"servesBeer"`
	ServesWine                   bool                 `json:"servesWine"`
	ServesBrunch                 bool                 `json:"servesBrunch"`
	ServesVegetarianFood         bool                 `json:"servesVegetarianFood"`
	OutdoorSeating               bool                 `json:"outdoorSeating"`
	LiveMusic                    bool                 `json:"liveMusic"`
	MenuForChildren              bool                 `json:"menuForChildren"`
	ServesCocktails              bool                 `json:"servesCocktails"`
	ServesDessert                bool                 `json:"servesDessert"`
	ServesCoffee                 bool                 `json:"servesCoffee"`
	GoodForChildren              bool                 `json:"goodForChildren"`
	AllowsDogs                   bool                 `json:"allowsDogs"`
	Restroom                     bool                 `json:"restroom"`
	GoodForGroups                bool                 `json:"goodForGroups"`
	GoodForWatchingSports        bool                 `json:"goodForWatchingSports"`
	AccessibilityOptions         AccessibilityOptions `json:"accessibilityOptions"`
	PureServiceAreaBusiness      bool                 `json:"pureServiceAreaBusiness"`
}
