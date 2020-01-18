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

type CityTestSuite struct {
	suite.Suite

	serverHandler http.Handler
	db            *dbx.DB
}

func (s *CityTestSuite) SetupTest() {
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

func (s *CityTestSuite) TestPatchCityOK() {
	city := entity.City{Name: "Munich", Latitude: 55.66, Longitude: 66.77}
	s.Require().NoError(s.db.Model(&city).Insert())

	newName := "NewMinuch"

	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPatch,
		fmt.Sprintf("/cities/%d", city.ID),
		[]byte(fmt.Sprintf(`{"name": "%s"}`, newName)),
	)

	s.Require().Equal(http.StatusOK, resp.Code)

	var b entity.City
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&b))

	s.Require().Equal(newName, b.Name)
}

func (s *CityTestSuite) TestDeleteCityOK() {
	city := entity.City{Name: "Ivanovo", Latitude: 55.66, Longitude: 66.77}
	require.NoError(s.T(), s.db.Model(&city).Insert())

	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodDelete,
		fmt.Sprintf("/cities/%d", city.ID),
		[]byte(nil),
	)

	require.Equal(s.T(), http.StatusOK, resp.Code)

	resp = runV1Request(s.T(),
		s.serverHandler,
		http.MethodDelete,
		fmt.Sprintf("/cities/%d", city.ID),
		[]byte(nil),
	)

	require.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *CityTestSuite) TestDeleteCityBadID() {
	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodDelete,
		fmt.Sprintf("/cities/%d", 0),
		[]byte(nil),
	)

	require.Equal(s.T(), http.StatusNotFound, resp.Code)
}
