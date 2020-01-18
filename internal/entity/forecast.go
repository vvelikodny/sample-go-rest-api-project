package entity

// Forecast represents an forecast for a particular city.
type Forecast struct {
	CityID int `json:"city_id"`
	Min    int `json:"min"`
	Max    int `json:"max"`
	Sample int `json:"sample"`
}
