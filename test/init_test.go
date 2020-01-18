package test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vvelikodny/weather/internal/config"
	"github.com/vvelikodny/weather/pkg/log"
)

type ValidationError struct {
	ErrorCode string `json:"error_code"`
	Details   []map[string]interface{}
}

func TestSuites(t *testing.T) {
	require.NoError(t, resetDB(t))

	suite.Run(t, new(CityTestSuite))
	suite.Run(t, new(TemperatureTestSuite))
	suite.Run(t, new(WebhookTestSuite))
}

func resetDB(t *testing.T) error {
	// load application configurations
	cfg, err := config.Load("../config/test.yml", log.New())
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", cfg.DSN)

	db.Query(`drop table if exists schema_migrations cascade`)
	db.Query(`drop table if exists temperature cascade`)
	db.Query(`drop table if exists webhook cascade`)
	db.Query(`drop table if exists city cascade`)

	runMigrations(db)

	return nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}

func runV1Request(t *testing.T, router http.Handler, method, URL string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(
		method,
		URL,
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
