package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"weatherservice"
)

const googlePhotosURL = "https://places.googleapis.com/v1/%s/media?maxHeightPx=400&maxWidthPx=400&key=%s"

type GooglePhotosAPI struct {
	client *http.Client

	key string
}

func NewGooglePhotosAPI(key string) *GooglePhotosAPI {
	return &GooglePhotosAPI{
		client: http.DefaultClient,
		key:    key,
	}
}

type hotelPhotoResult struct {
	idx       int
	photoUrls string
}

func (api *GooglePhotosAPI) FetchReportData(ctx context.Context, hotelName ...string) (*weatherservice.DataToReport[weatherservice.HotelsPhotos], error) {
	hotelPhotos := make(weatherservice.HotelsPhotos, len(hotelName))

	wg := sync.WaitGroup{}

	photoUrls := make(chan hotelPhotoResult, len(hotelName))

	for i, place := range hotelName {
		params := url.Values{}
		params.Add("key", api.key)

		// introduce grorutines with buffered channel and wait group
		wg.Add(1)
		go func(hotelName string, index int, photoUrlsch chan<- hotelPhotoResult) {
			defer wg.Done()
			if err := api.getImageFromGooglePlaces(hotelName, index, photoUrlsch); err != nil {
				log.Printf("Error fetching photo for %s: %v", hotelName, err)
			}
		}(place, i, photoUrls)
	}

	// goroutine to wait for all the responses and close the channel
	go func() {
		wg.Wait()
		close(photoUrls)
	}()

	// add the response of each goroutine to the hotelPhotos
	for img := range photoUrls {
		hotelPhotos[img.idx].HotelPhotos = append(hotelPhotos[img.idx].HotelPhotos, img.photoUrls)
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

func (api *GooglePhotosAPI) getImageFromGooglePlaces(place string, index int, out chan<- hotelPhotoResult) error {
	photoUrl := fmt.Sprintf(googlePhotosURL, place, api.key)

	req, err := http.NewRequest("GET", photoUrl, nil)
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	img := base64.StdEncoding.EncodeToString(body)

	out <- hotelPhotoResult{
		idx:       index,
		photoUrls: img,
	}

	return nil
}

func (api *GooglePhotosAPI) FetchGeneralInfo(ctx context.Context, _ ...string) (*weatherservice.DataToReport[weatherservice.HotelsPhotos], error) {
	return nil, nil
}

type SearchRequest struct {
	TextQuery string `json:"textQuery"`
}

type RoutingSummary struct {
	// Define fields for the RoutingSummary object as needed
}

type ContextualContent struct {
	// Define fields for the ContextualContent object as needed
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
