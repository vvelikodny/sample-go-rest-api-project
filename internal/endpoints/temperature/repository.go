package temperature

import (
	"context"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

// Repository encapsulates the logic to access temperatures from the data source.
type Repository interface {
	// Create saves a new temperature in the storage.
	Get(ctx context.Context, int int) (entity.Temperature, error)
	// Create saves a new temperature in the storage.
	Create(ctx context.Context, temperature *entity.Temperature) error
}

// repository persists temperatures in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new temperature repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id int) (entity.Temperature, error) {
	var temperature entity.Temperature
	err := r.db.With(ctx).Select().Model(id, &temperature)
	return temperature, err
}

// Create saves a new temperature record in the database.
// It returns the ID of the newly inserted temperature record.
func (r repository) Create(ctx context.Context, temperature *entity.Temperature) error {
	return r.db.With(ctx).Model(temperature).Insert()
}
