package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"backend/internal/handlers"
)

// test inserting a valid swift code
func TestPostSwiftCode_Success(t *testing.T) {
	t.Log("Testing successful SWIFT code insertion")

	handler := handlers.NewHandler(db)

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address",
		"isHeadquarter": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SWIFT code added successfully", response["message"])
}

// test inserting a duplicate swift code
func TestPostSwiftCode_Conflict(t *testing.T) {
	t.Log("Testing duplicate SWIFT code should return 409 Conflict")

	handler := handlers.NewHandler(db)

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address",
		"isHeadquarter": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"], "SWIFT code already exists")
}

// test inserting an invalid json request
func TestPostSwiftCode_InvalidJSON(t *testing.T) {
	t.Log("Testing invalid JSON should return 400 Bad Request")

	handler := handlers.NewHandler(db)

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid JSON")
}

// test inserting a swift code with missing required fields
func TestPostSwiftCode_MissingFields(t *testing.T) {
	t.Log("Testing missing required fields should return 400 Bad Request")

	handler := handlers.NewHandler(db)

	body := `{ "swiftCode": "ABCDEFGHXXX" }`

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Missing required fields")
}

// test handling database failure
func TestPostSwiftCode_DatabaseError(t *testing.T) {
	t.Log("Testing database error should return 500 Internal Server Error")

	handler := handlers.NewHandler(db)

	// forcefully drop connection
	db.Close()

	body := `{
		"swiftCode": "ERROR123XXX",
		"bankName": "Error Bank",
		"countryISO2": "XX",
		"countryName": "Errorland",
		"address": "Error Address",
		"isHeadquarter": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to insert SWIFT code")
}

// test inserting a swift code without isHeadquarter field
func TestPostSwiftCode_MissingIsHeadquarter(t *testing.T) {
	t.Log("Testing missing isHeadquarter field should return 400 Bad Request")

	handler := handlers.NewHandler(db)

	// missing "isHeadquarter" field in the request body
	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address"
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Missing required fields")
}
