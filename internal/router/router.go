package router

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/vvelikodny/weather/internal/config"
	"github.com/vvelikodny/weather/internal/endpoints/city"
	"github.com/vvelikodny/weather/internal/endpoints/forecast"
	"github.com/vvelikodny/weather/internal/endpoints/temperature"
	"github.com/vvelikodny/weather/internal/endpoints/webhook"
	"github.com/vvelikodny/weather/internal/errors"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func BuildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	rg := router.Group("/v1")

	cityRepo := city.NewRepository(db, logger)

	city.RegisterHandlers(rg,
		city.NewService(cityRepo, logger),
		logger,
	)

	temperature.RegisterHandlers(rg,
		temperature.NewService(temperature.NewRepository(db, logger), logger),
		logger,
	)

	forecast.RegisterHandlers(rg,
		forecast.NewService(forecast.NewRepository(db, logger), logger),
		logger,
	)

	webhook.RegisterHandlers(rg,
		webhook.NewService(webhook.NewRepository(db, logger, cityRepo), logger),
		logger,
	)

	return router
}
