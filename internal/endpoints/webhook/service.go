package webhook

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/pkg/log"
)

// Service encapsulates logic for webhooks.
type Service interface {
	Create(ctx context.Context, input CreateWebhookRequest) (Webhook, error)
	Delete(ctx context.Context, id int) (Webhook, error)
}

// Webhook represents the data about an webhook.
type Webhook struct {
	entity.Webhook
}

// CreateWebhookRequest represents an webhook creation request.
type CreateWebhookRequest struct {
	CityID      int    `json:"city_id"`
	CallbackURL string `json:"callback_url"`
}

// Validate validates the CreateWebhookRequest fields.
func (m CreateWebhookRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CityID, validation.Required),
		validation.Field(&m.CallbackURL, validation.Required, is.URL),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new webhook service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the webhook with the specified the webhook ID.
func (s service) Get(ctx context.Context, id int) (Webhook, error) {
	webhook, err := s.repo.Get(ctx, id)
	if err != nil {
		return Webhook{}, err
	}
	return Webhook{webhook}, nil
}

// Create creates a new webhook.
func (s service) Create(ctx context.Context, req CreateWebhookRequest) (Webhook, error) {
	if err := req.Validate(); err != nil {
		return Webhook{}, err
	}

	webhook := entity.Webhook{
		CityID:      req.CityID,
		CallbackURL: req.CallbackURL,
	}

	err := s.repo.Create(ctx, &webhook)
	if err != nil {
		return Webhook{}, err
	}
	return s.Get(ctx, webhook.ID)
}

// Delete deletes the webhook with the specified ID.
func (s service) Delete(ctx context.Context, id int) (Webhook, error) {
	city, err := s.Get(ctx, id)
	if err != nil {
		return Webhook{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Webhook{}, err
	}
	return city, nil
}
