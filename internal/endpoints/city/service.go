package city

import (
	"context"
	"reflect"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/log"
)

// Service encapsulates logic for cities.
type Service interface {
	Create(ctx context.Context, input CreateCityRequest) (City, error)
	Update(ctx context.Context, id int, input PatchCityRequest) (City, error)
	Delete(ctx context.Context, id int) (City, error)
}

// City represents the data about an city.
type City struct {
	entity.City
}

// CreateCityRequest represents an city creation request.
type CreateCityRequest struct {
	Name      string  `json:"name" `
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Validate validates the CreateCityRequest fields.
func (m CreateCityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 128)),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Name, validation.Required),
	)
}

// PatchCityRequest represents an city patch request.
type PatchCityRequest struct {
	Name      *string  `json:"name,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// Validate validates the CreateCityRequest fields.
func (m PatchCityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.NilOrNotEmpty, validation.Length(1, 128)),
		validation.Field(&m.Latitude, validation.NilOrNotEmpty),
		validation.Field(&m.Longitude, validation.NilOrNotEmpty),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new city service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the city with the specified the city ID.
func (s service) Get(ctx context.Context, id int) (City, error) {
	city, err := s.repo.Get(ctx, id)
	if err != nil {
		return City{}, err
	}
	return City{city}, nil
}

// Create creates a new city.
func (s service) Create(ctx context.Context, req CreateCityRequest) (City, error) {
	if err := req.Validate(); err != nil {
		return City{}, err
	}
	now := time.Now()
	city := entity.City{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: now,
	}
	err := s.repo.Create(ctx, &city)
	if err != nil {
		return City{}, err
	}
	return s.Get(ctx, city.ID)
}

// Update updates the city with the specified ID.
func (s service) Update(ctx context.Context, id int, req PatchCityRequest) (City, error) {
	if err := req.Validate(); err != nil {
		return City{}, err
	}

	city, err := s.Get(ctx, id)
	if err != nil {
		return city, err
	}

	if !patchValue(s.logger, &city, req) {
		return city, nil
	}

	if err := s.repo.Update(ctx, city.City); err != nil {
		return city, err
	}
	return city, nil
}

// Delete deletes the city with the specified ID.
func (s service) Delete(ctx context.Context, id int) (City, error) {
	city, err := s.Get(ctx, id)
	if err != nil {
		return City{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return City{}, err
	}
	return city, nil
}

func patchValue(logger log.Logger, entity interface{}, req PatchCityRequest) bool {
	rt := reflect.TypeOf(req)
	// reflect.Type
	rv := reflect.ValueOf(req)
	// reflect.Value
	cityv := reflect.ValueOf(entity)
	// reflect.Value
	patch := false
	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			cityv.Elem().FieldByName(rt.Field(i).Name).Set(rv.Field(i).Elem())

			patch = true
		}
	}
	return patch
}
