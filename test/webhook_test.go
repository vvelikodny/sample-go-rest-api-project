package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vvelikodny/weather/internal/entity"
	"net/http"
	"os"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/stretchr/testify/suite"
	"github.com/vvelikodny/weather/internal/config"
	"github.com/vvelikodny/weather/internal/router"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

type WebhookTestSuite struct {
	suite.Suite

	serverHandler http.Handler
	db            *dbx.DB
}

func (s *WebhookTestSuite) SetupTest() {
	logger := log.New()

	var err error
	// load application configurations
	cfg, err := config.Load("../config/test.yml", logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	s.db, err = dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	s.serverHandler = router.BuildHandler(logger, dbcontext.New(s.db), cfg)
}

func (s *WebhookTestSuite) TestCreateWebhookEmptyBody() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/webhooks",
		[]byte(``),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *WebhookTestSuite) TestCreateWebhookEmptyJSON() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/webhooks",
		[]byte(`{}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Len(b.Details, 2)
}

func (s *WebhookTestSuite) TestCreateWebhookEmptyURL() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/webhooks",
		[]byte(`{"city_id": 111}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "callback_url")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *WebhookTestSuite) TestCreateWebhookBadURL() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/webhooks",
		[]byte(`{"city_id": 111, "callback_url": "url"}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "callback_url")
	s.Contains(b.Details[0]["error"], "must be a valid URL")
}

func (s *WebhookTestSuite) TestCreateWebhookEmptyCityID() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/webhooks",
		[]byte(`{"callback_url": "https://www.finleap.com"}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "city_id")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *WebhookTestSuite) TestDeleteWebhookOK() {
	city := entity.City{Name: "Moscow", Latitude: 55.66, Longitude: 66.77}
	s.Require().NoError(s.db.Model(&city).Insert())
	webhook := entity.Webhook{CityID: city.ID, CallbackURL: "https://finleap.com"}
	s.Require().NoError(s.db.Model(&webhook).Insert())

	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodDelete,
		fmt.Sprintf("/webhooks/%d", webhook.ID),
		[]byte(nil),
	)

	require.Equal(s.T(), http.StatusOK, resp.Code)
}
