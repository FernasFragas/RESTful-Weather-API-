package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"weatherservice"
)

const foursquareSearchURL = "https://places-api.foursquare.com/places/search?"

type FoursquareAPI struct {
	client *http.Client

	key string
}

func NewFoursquareAPI(key string) *FoursquareAPI {
	return &FoursquareAPI{
		client: http.DefaultClient,
		key:    key,
	}
}

func (api *FoursquareAPI) FetchReportData(ctx context.Context, categories ...string) (*weatherservice.DataToReport[weatherservice.Coordinates], error) {
	if api.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	categoriesWithLonLat := strings.Split(categories[0], ",")
	categoriesWithoutLonLat := categoriesWithLonLat[2:]

	categoriesToSearchIDs := api.filterCategories(categoriesWithoutLonLat)

	queryParams := url.Values{}
	queryParams.Add("ll", fmt.Sprintf("%s,%s", categoriesWithLonLat[1], categoriesWithLonLat[0]))
	queryParams.Add("limit", "10")
	queryParams.Add("radius", "5000")
	queryParams.Add("categories", categoriesToSearchIDs)

	apiUrl := fmt.Sprintf("%s%s", foursquareSearchURL, queryParams.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", api.key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Places-Api-Version", "2025-06-17")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch data from geoapify: %s, because: %s", resp.Status, string(body))
	}

	fmt.Println(resp.Body)

	var foursquareResponse FoursquareResponse
	err = json.NewDecoder(resp.Body).Decode(&foursquareResponse)
	if err != nil {
		return nil, err
	}

	fmt.Println(foursquareResponse)

	// todo array with all the coorodinates of the places
	coordinates := make([]weatherservice.CoordinatesWithName, 0)
	for _, place := range foursquareResponse.Results {
		coordinates = append(coordinates, weatherservice.CoordinatesWithName{
			Name:        *place.Name,
			Coordinates: []interface{}{*place.Latitude, *place.Longitude},
		})
	}

	return &weatherservice.DataToReport[weatherservice.Coordinates]{
		Data: weatherservice.Coordinates{
			CoordinatesWithName: coordinates,
		},
	}, nil

}

func (api *FoursquareAPI) FetchGeneralInfo(ctx context.Context, city ...string) (*weatherservice.DataToReport[weatherservice.Coordinates], error) {
	return nil, nil
}

func (api *FoursquareAPI) filterCategories(categories []string) string {
	filteredCategories := []string{}

	for _, category := range categories {
		// Get individual features for this category
		features := api.availableCategories(category)
		filteredCategories = append(filteredCategories, features...)
	}

	return strings.Join(filteredCategories, ",")
}

func (api *FoursquareAPI) availableCategories(category string) []string {
	category = strings.ToLower(category)
	switch category {
	case "restaurants":
		return []string{
			"13065", // Restaurant
			"13003", // Bar
			"13002", // Fast Food
			"13032", // Coffee Shop
		}
	case "tourism":
		return []string{
			"16000", // Landmarks and Outdoors
			"12000", // Arts and Entertainment
			"10000", // Events
		}
	case "entertainment":
		return []string{
			"12000", // Arts and Entertainment
			"12001", // Aquarium
			"12002", // Arcade
			"12003", // Bowling Alley
			"12004", // Casino
			"12005", // Cinema
		}
	case "hotels":
		return []string{
			"19014", // Hotel
			"19015", // Hostel
			"19013", // Bed and Breakfast
		}
	case "shopping":
		return []string{
			"17000", // Retail
			"17001", // Adult Boutique
			"17002", // Antique Store
			"17003", // Arts and Crafts Store
		}
	case "outdoor":
		return []string{
			"16000", // Landmarks and Outdoors
			"16001", // Beach
			"16002", // Bridge
			"16003", // Cemetery
			"16004", // Farm
			"16005", // Garden
			"16006", // Lake
			"16007", // Mountain
			"16008", // National Park
		}
	default:
		return []string{}
	}
}

// ... existing code ...

// Place represents a place from Foursquare Places API
type FoursquarePlace struct {
	FSQPlaceId       string            `json:"fsq_place_id"`
	Latitude         *float64          `json:"latitude,omitempty"`
	Longitude        *float64          `json:"longitude,omitempty"`
	Categories       []Category        `json:"categories,omitempty"`
	Chains           []Chain           `json:"chains,omitempty"`
	DateClosed       *string           `json:"date_closed,omitempty"`
	DateCreated      *string           `json:"date_created,omitempty"`
	DateRefreshed    *string           `json:"date_refreshed,omitempty"`
	Description      *string           `json:"description,omitempty"`
	Distance         *int              `json:"distance,omitempty"`
	Email            *string           `json:"email,omitempty"`
	ExtendedLocation *ExtendedLocation `json:"extended_location,omitempty"`
	Attributes       *Attributes       `json:"attributes,omitempty"`
	Hours            *Hours            `json:"hours,omitempty"`
	HoursPopular     []HourPeriod      `json:"hours_popular,omitempty"`
	Link             *string           `json:"link,omitempty"`
	Location         *Location         `json:"location,omitempty"`
	Menu             *string           `json:"menu,omitempty"`
	Name             *string           `json:"name,omitempty"`
	Photos           []Photo           `json:"photos,omitempty"`
	Popularity       *int              `json:"popularity,omitempty"`
	PlacemakerURL    *string           `json:"placemaker_url,omitempty"`
	Price            *int              `json:"price,omitempty"`
	Rating           *float64          `json:"rating,omitempty"`
	RelatedPlaces    *RelatedPlaces    `json:"related_places,omitempty"`
	SocialMedia      *SocialMedia      `json:"social_media,omitempty"`
	Stats            *Stats            `json:"stats,omitempty"`
	StoreID          *string           `json:"store_id,omitempty"`
	Tastes           []string          `json:"tastes,omitempty"`
	Tel              *string           `json:"tel,omitempty"`
	Tips             []Tip             `json:"tips,omitempty"`
	Verified         *bool             `json:"verified,omitempty"`
	Website          *string           `json:"website,omitempty"`
}

type Category struct {
	FSQCategoryID string `json:"fsq_category_id"`
	Name          string `json:"name"`
	ShortName     string `json:"short_name"`
	PluralName    string `json:"plural_name"`
	Icon          *Icon  `json:"icon,omitempty"`
}

type Icon struct {
	ID              string   `json:"id"`
	CreatedAt       string   `json:"created_at"`
	Prefix          string   `json:"prefix"`
	Suffix          string   `json:"suffix"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	Classifications []string `json:"classifications,omitempty"`
	Tip             *Tip     `json:"tip,omitempty"`
}

type Chain struct {
	FSQChainID string `json:"fsq_chain_id"`
	Name       string `json:"name"`
	Logo       *Icon  `json:"logo,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
}

type ExtendedLocation struct {
	DMA           string `json:"dma,omitempty"`
	CensusBlockID string `json:"census_block_id,omitempty"`
}

type Attributes struct {
	Restroom        interface{} `json:"restroom,omitempty"`
	OutdoorSeating  interface{} `json:"outdoor_seating,omitempty"`
	ATM             interface{} `json:"atm,omitempty"`
	HasParking      interface{} `json:"has_parking,omitempty"`
	WiFi            *string     `json:"wifi,omitempty"`
	Delivery        interface{} `json:"delivery,omitempty"`
	Reservations    interface{} `json:"reservations,omitempty"`
	TakesCreditCard interface{} `json:"takes_credit_card,omitempty"`
}

type Location struct {
	Address          string `json:"address,omitempty"`
	Locality         string `json:"locality,omitempty"`
	Region           string `json:"region,omitempty"`
	Postcode         string `json:"postcode,omitempty"`
	AdminRegion      string `json:"admin_region,omitempty"`
	PostTown         string `json:"post_town,omitempty"`
	POBox            string `json:"po_box,omitempty"`
	Country          string `json:"country,omitempty"`
	FormattedAddress string `json:"formatted_address,omitempty"`
}

type Hours struct {
	Display        string       `json:"display,omitempty"`
	IsLocalHoliday *bool        `json:"is_local_holiday,omitempty"`
	OpenNow        *bool        `json:"open_now,omitempty"`
	Regular        []HourPeriod `json:"regular,omitempty"`
}

type HourPeriod struct {
	Close string `json:"close"`
	Day   int    `json:"day"`
	Open  string `json:"open"`
}

type FoursquarePhoto struct {
	FSQPhotoID      string   `json:"fsq_photo_id,omitempty"`
	ID              string   `json:"id,omitempty"`
	CreatedAt       string   `json:"created_at"`
	Prefix          string   `json:"prefix"`
	Suffix          string   `json:"suffix"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	Classifications []string `json:"classifications,omitempty"`
	Tip             *Tip     `json:"tip,omitempty"`
}

type Tip struct {
	FSQTipID      string           `json:"fsq_tip_id,omitempty"`
	ID            string           `json:"id,omitempty"`
	CreatedAt     string           `json:"created_at"`
	Text          string           `json:"text"`
	URL           string           `json:"url,omitempty"`
	Photo         *FoursquarePhoto `json:"photo,omitempty"`
	Lang          string           `json:"lang,omitempty"`
	AgreeCount    int              `json:"agree_count"`
	DisagreeCount int              `json:"disagree_count"`
}

type RelatedPlaces struct {
	Parent   *Place  `json:"parent,omitempty"`
	Children []Place `json:"children,omitempty"`
}

type SocialMedia struct {
	FacebookID string `json:"facebook_id,omitempty"`
	Instagram  string `json:"instagram,omitempty"`
	Twitter    string `json:"twitter,omitempty"`
}

type Stats struct {
	TotalPhotos  int `json:"total_photos"`
	TotalRatings int `json:"total_ratings"`
	TotalTips    int `json:"total_tips"`
}

type Context struct {
	GeoBounds *GeoBounds `json:"geo_bounds,omitempty"`
}

type GeoBounds struct {
	Circle *Circle `json:"circle,omitempty"`
}

type Circle struct {
	Center *Center `json:"center,omitempty"`
	Radius int     `json:"radius,omitempty"`
}

type Center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// API Response structures
type FoursquareResponse struct {
	Results []FoursquarePlace `json:"results"`
	Context *Context          `json:"context,omitempty"`
}
