package city

import (
	"context"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

// Repository encapsulates the logic to access cities from the data source.
type Repository interface {
	// Create saves a new city in the storage.
	Get(ctx context.Context, int int) (entity.City, error)
	// Create saves a new city in the storage.
	Create(ctx context.Context, city *entity.City) error
	// Update updates the city with given ID in the storage.
	Update(ctx context.Context, city entity.City) error
	// Delete removes the city with given ID from the storage.
	Delete(ctx context.Context, id int) error
}

// repository persists cities in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new city repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id int) (entity.City, error) {
	var city entity.City
	err := r.db.With(ctx).Select().Model(id, &city)
	return city, err
}

// Create saves a new city record in the database.
// It returns the ID of the newly inserted city record.
func (r repository) Create(ctx context.Context, city *entity.City) error {
	return r.db.With(ctx).Model(city).Insert()
}

// Update saves the changes to an city in the database.
func (r repository) Update(ctx context.Context, city entity.City) error {
	return r.db.With(ctx).Model(&city).Update()
}

// Delete deletes an city with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id int) error {
	city, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&city).Delete()
}
