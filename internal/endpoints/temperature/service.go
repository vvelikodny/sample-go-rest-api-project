package temperature

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/log"
)

// Service encapsulates logic for temperature.
type Service interface {
	Create(ctx context.Context, input CreateTemperatureRequest) (Temperature, error)
}

// Temperature represents the data about an temperature.
type Temperature struct {
	entity.Temperature
}

// CreateTemperatureRequest represents an temperature creation request.
type CreateTemperatureRequest struct {
	CityID int  `json:"city_id"`
	Min    *int `json:"min"`
	Max    *int `json:"max"`
}

// Validate validates the CreateTemperatureRequest fields.
func (m CreateTemperatureRequest) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.CityID, validation.Required),
		validation.Field(&m.Min, validation.Required, validation.Min(-100), validation.Max(100)),
		validation.Field(&m.Max, validation.Required, validation.Min(-100), validation.Max(100)),
	)
	if err != nil {
		return err
	}

	errs := validation.Errors{}
	if *m.Min > *m.Max {
		errs["min"] = errors.New("min should be less then max")
		return errs
	}

	return nil
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new temperature service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the temperature with the specified the temperature ID.
func (s service) Get(ctx context.Context, id int) (Temperature, error) {
	temperature, err := s.repo.Get(ctx, id)
	if err != nil {
		return Temperature{}, err
	}
	return Temperature{temperature}, nil
}

// Create creates a new temperature.
func (s service) Create(ctx context.Context, req CreateTemperatureRequest) (Temperature, error) {
	if err := req.Validate(); err != nil {
		return Temperature{}, err
	}
	now := time.Now()
	temperature := entity.Temperature{
		CityID:    req.CityID,
		Min:       *req.Min,
		Max:       *req.Max,
		CreatedAt: now,
	}
	err := s.repo.Create(ctx, &temperature)
	if err != nil {
		return Temperature{}, err
	}
	return s.Get(ctx, temperature.ID)
}
