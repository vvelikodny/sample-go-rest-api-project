package webhook

import (
	"net/http"
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/vvelikodny/weather/internal/errors"
	"github.com/vvelikodny/weather/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/webhooks", res.create)
	r.Delete("/webhooks/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) create(c *routing.Context) error {
	var input CreateWebhookRequest
	if err := c.Read(&input); err != nil {
		return errors.BadRequest("")
	}
	webhook, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(webhook, http.StatusCreated)
}

func (r resource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.BadRequest("")
	}

	city, err := r.service.Delete(c.Request.Context(), id)
	if err != nil {
		return err
	}

	return c.Write(city)
}
