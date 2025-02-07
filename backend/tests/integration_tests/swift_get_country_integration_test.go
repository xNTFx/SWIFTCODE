package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"backend/internal/handlers"
	"backend/internal/models"
)

// test retrieving swift codes by country (valid country)
func TestGetSwiftCodesByCountry_Success(t *testing.T) {
	t.Log("Testing retrieval of SWIFT codes for a valid country (AL)")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/AL", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.SwiftCodeByCountryISO2
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "AL", response.CountryISO2)
	assert.Equal(t, "ALBANIA", response.CountryName)
	assert.Greater(t, len(response.SwiftCodes), 0) // should return at least 1 swift code
}

// test retrieving swift codes by country (non-existent country)
func TestGetSwiftCodesByCountry_NotFound(t *testing.T) {
	t.Log("Testing retrieval of SWIFT codes for a non-existent country should return 404 Not Found")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/XX", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// test retrieving swift codes by country with an invalid format
func TestGetSwiftCodesByCountry_InvalidFormat(t *testing.T) {
	t.Log("Testing retrieval of SWIFT codes with an invalid ISO2 code format should return 400 Bad Request")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/INVALID!@#", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
