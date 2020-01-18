package forecast

import (
	"context"
	"fmt"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/log"
)

// Service encapsulates logic for temperature.
type Service interface {
	Get(ctx context.Context, cityID int) (Forecast, error)
}

// Forecast represents the data about an forecast.
type Forecast struct {
	entity.Forecast
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
func (s service) Get(ctx context.Context, id int) (Forecast, error) {
	forecast, err := s.repo.Get(ctx, id)
	if err != nil {
		return Forecast{}, fmt.Errorf("could'n get forecast from db %w", err)
	}
	return Forecast{forecast}, nil
}
