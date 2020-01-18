package temperature

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/vvelikodny/weather/internal/errors"
	"github.com/vvelikodny/weather/pkg/log"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/temperatures", res.create)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) create(c *routing.Context) error {
	var input CreateTemperatureRequest
	if err := c.Read(&input); err != nil {
		return errors.BadRequest("")
	}
	temperature, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(temperature, http.StatusCreated)
}
