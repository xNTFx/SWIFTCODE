package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handlers"
	"backend/internal/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// test retrieving an existing headquarters SWIFT code
func TestGetSwiftCodeDetails_Headquarter(t *testing.T) {
	t.Log("Testing retrieval of an existing headquarters SWIFT code")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/AAISALTRXXX", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodeDetailsHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.SwiftCodeHeadquarter
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "AAISALTRXXX", response.SwiftCode)
	assert.Equal(t, "UNITED BANK OF ALBANIA SH.A", response.BankName)
	assert.Equal(t, "AL", response.CountryISO2)
	assert.Equal(t, "ALBANIA", response.CountryName)
	assert.Equal(t, true, response.IsHeadquarter)
}

// test retrieving a non-existent SWIFT code
func TestGetSwiftCodeDetails_NotFound(t *testing.T) {
	t.Log("Testing retrieval of a non-existent SWIFT code should return 404 Not Found")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/NONEXISTENT", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodeDetailsHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// test retrieving a SWIFT code with an invalid format
func TestGetSwiftCodeDetails_InvalidFormat(t *testing.T) {
	t.Log("Testing retrieval of an invalid SWIFT code format should return 400 Bad Request")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/INVALID_CODE!", nil)
	rec := httptest.NewRecorder()

	handler.GetSwiftCodeDetailsHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
