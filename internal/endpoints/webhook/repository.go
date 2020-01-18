package webhook

import (
	"context"
	"fmt"

	"github.com/vvelikodny/weather/internal/endpoints/city"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

// Repository encapsulates the logic to access webhooks from the data source.
type Repository interface {
	// Create saves a new webhook in the storage.
	Get(ctx context.Context, int int) (entity.Webhook, error)
	// Create saves a new webhook in the storage.
	Create(ctx context.Context, webhook *entity.Webhook) error
	// Delete removes the webhook with given ID from the storage.
	Delete(ctx context.Context, id int) error
}

// repository persists webhooks in database
type repository struct {
	db             *dbcontext.DB
	logger         log.Logger
	cityRepository city.Repository
}

// NewRepository creates a new webhook repository
func NewRepository(db *dbcontext.DB, logger log.Logger, cityRepository city.Repository) Repository {
	return repository{db, logger, cityRepository}
}

func (r repository) Get(ctx context.Context, id int) (entity.Webhook, error) {
	var webhook entity.Webhook
	err := r.db.With(ctx).Select().Model(id, &webhook)
	return webhook, err
}

// Create saves a new webhook record in the database.
// It returns the ID of the newly inserted webhook record.
func (r repository) Create(ctx context.Context, webhook *entity.Webhook) error {
	_, err := r.cityRepository.Get(ctx, webhook.CityID)
	if err != nil {
		return fmt.Errorf("city %v: %w", webhook.CityID, err)
	}

	return r.db.With(ctx).Model(webhook).Insert()
}

// Delete deletes an webhook with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id int) error {
	_, err := r.cityRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	webhook, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&webhook).Delete()
}
