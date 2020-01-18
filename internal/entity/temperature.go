package entity

import (
	"time"
)

// Temperature represents an temperature record.
type Temperature struct {
	ID        int       `json:"id"`
	CityID    int       `json:"city_id"`
	Min       int       `json:"min"`
	Max       int       `json:"max"`
	CreatedAt time.Time `json:"timestamp"`
}
