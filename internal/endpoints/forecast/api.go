package forecast

import (
	"fmt"
	"net/http"
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/vvelikodny/weather/internal/errors"
	"github.com/vvelikodny/weather/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/forecasts/<city_id>", res.get)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	cityId, err := strconv.Atoi(c.Param("city_id"))
	if err != nil {
		return errors.BadRequest("")
	}

	forecast, err := r.service.Get(c.Request.Context(), cityId)
	if err != nil {
		return fmt.Errorf("call forecast service %w", err)
	}

	return c.WriteWithStatus(forecast, http.StatusOK)
}
