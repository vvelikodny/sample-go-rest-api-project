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

type TemperatureTestSuite struct {
	suite.Suite

	serverHandler http.Handler
	db            *dbx.DB
}

func (s *TemperatureTestSuite) SetupTest() {
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

func (s *TemperatureTestSuite) TestCreateTemperatureEmptyBody() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(``),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *TemperatureTestSuite) TestCreateTemperatureEmptyJSON() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(`{}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "city_id")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *TemperatureTestSuite) TestCreateTemperatureBadMin() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(`{"city_id": 111, "max": 5}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "min")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *TemperatureTestSuite) TestCreateTemperatureBadMax() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(`{"city_id": 222, "min": 1}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "max")
	s.Contains(b.Details[0]["error"], "cannot be blank")
}

func (s *TemperatureTestSuite) TestCreateTemperatureMinLessThenMax() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(`{"city_id": 222, "min": 5, "max": 1}`),
	)

	require.Equal(s.T(), http.StatusBadRequest, resp.Code)

	var b ValidationError
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&b))

	s.NotNil(b.Details)
	s.Contains(b.Details[0], "field")
	s.Contains(b.Details[0]["field"], "min")
	s.Contains(b.Details[0]["error"], "min should be less then max")
}

func (s *TemperatureTestSuite) TestCreateTemperatureOK() {
	city := entity.City{Name: "Munich", Latitude: 55.66, Longitude: 66.77}
	s.Require().NoError(s.db.Model(&city).Insert())

	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(fmt.Sprintf(`{"city_id": %d, "min": 1, "max": 2}`, city.ID)),
	)

	require.Equal(s.T(), http.StatusCreated, resp.Code)
}
