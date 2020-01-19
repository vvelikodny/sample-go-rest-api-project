package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/stretchr/testify/suite"
	"github.com/vvelikodny/weather/internal/config"
	"github.com/vvelikodny/weather/internal/entity"
	"github.com/vvelikodny/weather/internal/router"
	"github.com/vvelikodny/weather/pkg/dbcontext"
	"github.com/vvelikodny/weather/pkg/log"
)

type ForecastTestSuite struct {
	suite.Suite

	serverHandler http.Handler
	db            *dbx.DB
}

func (s *ForecastTestSuite) SetupTest() {
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

func (s *TemperatureTestSuite) TestGetForecastOK() {
	city := entity.City{Name: "Munich", Latitude: 55.66, Longitude: 66.77}
	s.Require().NoError(s.db.Model(&city).Insert())

	resp := runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(fmt.Sprintf(`{"city_id": %d, "min": 1, "max": 2}`, city.ID)),
	)
	s.Require().Equal(http.StatusCreated, resp.Code)

	resp = runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(fmt.Sprintf(`{"city_id": %d, "min": -11, "max": 5}`, city.ID)),
	)
	s.Require().Equal(http.StatusCreated, resp.Code)

	resp = runV1Request(s.T(),
		s.serverHandler,
		http.MethodPost,
		"/temperatures",
		[]byte(fmt.Sprintf(`{"city_id": %d, "min": 4, "max": 15}`, city.ID)),
	)
	s.Require().Equal(http.StatusCreated, resp.Code)

	resp = runV1Request(s.T(),
		s.serverHandler,
		http.MethodGet,
		fmt.Sprintf("/forecasts/%d", city.ID),
		[]byte(nil),
	)
	s.Require().Equal(http.StatusOK, resp.Code)

	var b entity.Forecast
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&b))

	s.Require().Equal(city.ID, b.CityID)
	s.Require().Equal(-11, b.Min)
	s.Require().Equal(15, b.Max)
	s.Require().Equal(3, b.Sample)

}
