package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ValidationError struct {
	ErrorCode string `json:"error_code"`
	Details   []map[string]interface{}
}

func TestSuites(t *testing.T) {
	suite.Run(t, new(CityTestSuite))
}

func runV1Request(t *testing.T, router http.Handler, method, URL string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("/v1%s", URL),
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
