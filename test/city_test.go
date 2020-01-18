package test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/stretchr/testify/suite"
	"github.com/vvelikodny/weather/internal/config"
	"github.com/vvelikodny/weather/internal/router"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

type CityTestSuite struct {
	suite.Suite

	serverHandler http.Handler
}

func (s *CityTestSuite) SetupTest() {
	logger := log.New()

	// load application configurations
	cfg, err := config.Load("../config/test.yml", logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	s.serverHandler = router.BuildHandler(logger, dbcontext.New(db), cfg)
}

func (s *CityTestSuite) TestCreateCityEmptyBody() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/cities",
		[]byte(``),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *CityTestSuite) TestCreateCityEmptyJSON() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/cities",
		[]byte(`{}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "name")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *CityTestSuite) TestCreateCityOK() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/cities",
		[]byte(`{"name": "Berlin", "latitude": 55.66, "longitude": 66.77}`),
	)

	require.Equal(s.T(), http.StatusCreated, resp.Code)
}
