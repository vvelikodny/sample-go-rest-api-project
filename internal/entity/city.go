package entity

import (
	"time"
)

// City represents an city record.
type City struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" sql:"name"`
	Latitude  float64   `json:"latitude" sql:"latitude"`
	Longitude float64   `json:"longitude" sql:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}
