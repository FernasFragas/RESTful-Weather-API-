package myApi

type WeatherData struct {
	Weather     Weather     `json:"weather"`
	Temperature Temperature `json:"main"`
	Wind        Wind        `json:"wind"`
	Sys         Sys         `json:"sys"`
	Name        string      `json:"name"`
}

type Weather []struct {
	WeatherState string `json:"main"`
	Description  string `json:"description"`
}

type Temperature struct {
	FeelsLike float64 `json:"feels_like"`
	Temp_min  float64 `json:"temp_min"`
	Temp_max  float64 `json:"tem_max"`
	Humidity  float64 `json:"humidity"`
}

type Sys struct {
	Country string  `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset  float64 `json:"sunset"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}
