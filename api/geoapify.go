package api

const geoapifySearchURL = "https://api.geoapify.com/v2/place-details?features=walk_10,walk_10.restaurant,radius_500,radius_500.restaurant"

type GeoapifyResponse struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties Properties             `json:"properties"`
	Bbox       []float64              `json:"bbox,omitempty"`
	Center     []float64              `json:"center,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
}

type Properties struct {
	Datasource    Datasource `json:"datasource"`
	Country       string     `json:"country"`
	CountryCode   string     `json:"country_code"`
	State         string     `json:"state"`
	County        string     `json:"county"`
	City          string     `json:"city"`
	Postcode      string     `json:"postcode"`
	District      string     `json:"district"`
	Suburb        string     `json:"suburb"`
	Street        string     `json:"street"`
	Housenumber   string     `json:"housenumber"`
	Lon           float64    `json:"lon"`
	Lat           float64    `json:"lat"`
	Distance      float64    `json:"distance"`
	ResultType    string     `json:"result_type"`
	Formatted     string     `json:"formatted"`
	AddressLine1  string     `json:"address_line1"`
	AddressLine2  string     `json:"address_line2"`
	Category      string     `json:"category"`
	Timezone      Timezone   `json:"timezone"`
	PlusCode      string     `json:"plus_code"`
	PlusCodeShort string     `json:"plus_code_short"`
	Rank          Rank       `json:"rank"`
	PlaceID       string     `json:"place_id"`
	Bbox          []float64  `json:"bbox"`
	StreetNumber  string     `json:"street_number"`
	HouseNumber   string     `json:"house_number"`
	Road          string     `json:"road"`
	Neighbourhood string     `json:"neighbourhood"`
	Quarter       string     `json:"quarter"`
	Hamlet        string     `json:"hamlet"`
	Village       string     `json:"village"`
	Town          string     `json:"town"`
	Municipality  string     `json:"municipality"`
	CityDistrict  string     `json:"city_district"`
	StateDistrict string     `json:"state_district"`
	ISO31662      string     `json:"ISO3166-2"`
	StateCode     string     `json:"state_code"`
}

type Datasource struct {
	Sourcename  string `json:"sourcename"`
	Attribution string `json:"attribution"`
	License     string `json:"license"`
	URL         string `json:"url"`
	Raw         Raw    `json:"raw"`
}

type Raw struct {
	Name          string `json:"name"`
	Street        string `json:"street"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
	Postcode      string `json:"postcode"`
	District      string `json:"district"`
	Suburb        string `json:"suburb"`
	Housenumber   string `json:"housenumber"`
	Lon           string `json:"lon"`
	Lat           string `json:"lat"`
	Formatted     string `json:"formatted"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	Category      string `json:"category"`
	ResultType    string `json:"result_type"`
	Rank          Rank   `json:"rank"`
	PlaceID       string `json:"place_id"`
	Bbox          Bbox   `json:"bbox"`
	StreetNumber  string `json:"street_number"`
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Neighbourhood string `json:"neighbourhood"`
	Quarter       string `json:"quarter"`
	Hamlet        string `json:"hamlet"`
	Village       string `json:"village"`
	Town          string `json:"town"`
	Municipality  string `json:"municipality"`
	CityDistrict  string `json:"city_district"`
	StateDistrict string `json:"state_district"`
	ISO31662      string `json:"ISO3166-2"`
	StateCode     string `json:"state_code"`
}

type Timezone struct {
	Name             string `json:"name"`
	OffsetSTD        string `json:"offset_STD"`
	OffsetSTDSeconds int    `json:"offset_STD_seconds"`
	OffsetDST        string `json:"offset_DST"`
	OffsetDSTSeconds int    `json:"offset_DST_seconds"`
	AbbreviationSTD  string `json:"abbreviation_STD"`
	AbbreviationDST  string `json:"abbreviation_DST"`
}

type Rank struct {
	Importance          float64 `json:"importance"`
	Popularity          float64 `json:"popularity"`
	Confidence          float64 `json:"confidence"`
	ConfidenceCityLevel float64 `json:"confidence_city_level"`
	MatchType           string  `json:"match_type"`
}

type Bbox struct {
	Lon1 float64 `json:"lon1"`
	Lat1 float64 `json:"lat1"`
	Lon2 float64 `json:"lon2"`
	Lat2 float64 `json:"lat2"`
}
