package forecast

import (
	"context"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"time"

	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

// Repository encapsulates the logic to access forecasts from the data source.
type Repository interface {
	// Create saves a new temperature in the storage.
	Get(ctx context.Context, int int) (entity.Forecast, error)
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

func (r repository) Get(ctx context.Context, cityId int) (entity.Forecast, error) {
	var forecast entity.Forecast
	err := r.db.With(ctx).
		NewQuery(fmt.Sprintf(`
          SELECT
            city_id, MIN(min), MAX(max), count(*) as sample
          FROM
            temperature
         WHERE
           city_id = {:city_id} AND created_at >= {:day_before}
         GROUP BY
           city_id
		`)).
		Bind(dbx.Params{"city_id": cityId, "day_before": time.Now().AddDate(0, 0, -1)}).
		One(&forecast)
	return forecast, err
}
