package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/status", NewWebStatusHandler().Get)

	req, err := http.NewRequest("GET", "/status", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var status *Status
	err = json.Unmarshal(rr.Body.Bytes(), &status)
	assert.NoError(t, err)

	assert.IsType(t, status, &Status{})
	assert.Equal(t, status.Status, "UP")
}
