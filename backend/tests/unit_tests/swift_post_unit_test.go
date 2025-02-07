package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"backend/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestPostSwiftCodeHandler_InvalidJSON verifies that an invalid JSON request returns a bad request error.
func TestPostSwiftCodeHandler_InvalidJSON(t *testing.T) {
	t.Log("Testing invalid JSON returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader("{invalid}"))
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON")
}

// TestPostSwiftCodeHandler_Success verifies that a valid SWIFT code is added successfully.
func TestPostSwiftCodeHandler_Success(t *testing.T) {
	t.Log("Testing successful SWIFT code insertion")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		WITH country_ins AS (
			INSERT INTO countries (iso2_code, name)
			VALUES ($1, $2)
			ON CONFLICT (iso2_code) DO NOTHING
			RETURNING id
		), country_sel AS (
			SELECT id FROM countries WHERE iso2_code = $1
			UNION ALL
			SELECT id FROM country_ins LIMIT 1
		), bank_ins AS (
			INSERT INTO banks (name, country_id)
			SELECT $3, id FROM country_sel
			ON CONFLICT (name, country_id) DO NOTHING
			RETURNING id
		), bank_sel AS (
			SELECT id FROM banks WHERE name = $3 AND country_id = (SELECT id FROM country_sel)
			UNION ALL
			SELECT id FROM bank_ins LIMIT 1
		), swift_ins AS (
			INSERT INTO swift_codes (swift_code, bank_id, is_headquarter, address)
			SELECT $4, (SELECT id FROM bank_sel), $5, $6
			ON CONFLICT (swift_code) DO NOTHING
			RETURNING swift_code
		)
		SELECT COUNT(*) FROM swift_ins;
	`)).WithArgs("PL", "POLAND", "TEST BANK", "ABCDEFGHXXX", true, "TEST ADDRESS").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address",
		"isHeadquarter": true
	}`

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"message":"SWIFT code added successfully"}`, w.Body.String())
}

// TestPostSwiftCodeHandler_Conflict verifies that inserting a duplicate SWIFT code returns a conflict error.
func TestPostSwiftCodeHandler_Conflict(t *testing.T) {
	t.Log("Testing duplicate SWIFT code returns 409 Conflict")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		WITH country_ins AS (
			INSERT INTO countries (iso2_code, name)
			VALUES ($1, $2)
			ON CONFLICT (iso2_code) DO NOTHING
			RETURNING id
		), country_sel AS (
			SELECT id FROM countries WHERE iso2_code = $1
			UNION ALL
			SELECT id FROM country_ins LIMIT 1
		), bank_ins AS (
			INSERT INTO banks (name, country_id)
			SELECT $3, id FROM country_sel
			ON CONFLICT (name, country_id) DO NOTHING
			RETURNING id
		), bank_sel AS (
			SELECT id FROM banks WHERE name = $3 AND country_id = (SELECT id FROM country_sel)
			UNION ALL
			SELECT id FROM bank_ins LIMIT 1
		), swift_ins AS (
			INSERT INTO swift_codes (swift_code, bank_id, is_headquarter, address)
			SELECT $4, (SELECT id FROM bank_sel), $5, $6
			ON CONFLICT (swift_code) DO NOTHING
			RETURNING swift_code
		)
		SELECT COUNT(*) FROM swift_ins;
	`)).WithArgs("PL", "POLAND", "TEST BANK", "ABCDEFGHXXX", true, "TEST ADDRESS").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address",
		"isHeadquarter": true
	}`

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.JSONEq(t, `{"error":"SWIFT code already exists"}`, w.Body.String())
}

// TestPostSwiftCodeHandler_MissingFields verifies that missing required fields return a bad request error.
func TestPostSwiftCodeHandler_MissingFields(t *testing.T) {
	t.Log("Testing missing required fields returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	body := `{ "swiftCode": "ABCDEFGHXXX" }`
	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing required fields")
}

// TestPostSwiftCodeHandler_InvalidData verifies that invalid field values return a bad request error.
func TestPostSwiftCodeHandler_InvalidData(t *testing.T) {
	t.Log("Testing invalid field values return 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	body := `{
		"swiftCode": "INVALID123",
		"bankName": "@Invalid Bank Name!",
		"countryISO2": "PLX",
		"countryName": "Pol4nd",
		"address": "Invalid Address @!",
		"isHeadquarter": true
	}`

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Validation errors")
}

// TestPostSwiftCodeHandler_DatabaseError verifies that a database error results in a 500 response.
func TestPostSwiftCodeHandler_DatabaseError(t *testing.T) {
	t.Log("Testing database error returns 500 Internal Server Error")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		WITH country_ins AS (
			INSERT INTO countries (iso2_code, name)
			VALUES ($1, $2)
			ON CONFLICT (iso2_code) DO NOTHING
			RETURNING id
		), country_sel AS (
			SELECT id FROM countries WHERE iso2_code = $1
			UNION ALL
			SELECT id FROM country_ins LIMIT 1
		), bank_ins AS (
			INSERT INTO banks (name, country_id)
			SELECT $3, id FROM country_sel
			ON CONFLICT (name, country_id) DO NOTHING
			RETURNING id
		), bank_sel AS (
			SELECT id FROM banks WHERE name = $3 AND country_id = (SELECT id FROM country_sel)
			UNION ALL
			SELECT id FROM bank_ins LIMIT 1
		), swift_ins AS (
			INSERT INTO swift_codes (swift_code, bank_id, is_headquarter, address)
			SELECT $4, (SELECT id FROM bank_sel), $5, $6
			ON CONFLICT (swift_code) DO NOTHING
			RETURNING swift_code
		)
		SELECT COUNT(*) FROM swift_ins;
	`)).WithArgs("PL", "POLAND", "TEST BANK", "ABCDEFGHXXX", true, "TEST ADDRESS").WillReturnError(sql.ErrConnDone)

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address",
		"isHeadquarter": true
	}`

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to insert SWIFT code")
}

// TestPostSwiftCodeHandler_MissingIsHeadquarter verifies that missing isHeadquarter field returns a bad request error.
func TestPostSwiftCodeHandler_MissingIsHeadquarter(t *testing.T) {
	t.Log("Testing missing isHeadquarter field returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	body := `{
		"swiftCode": "ABCDEFGHXXX",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "Poland",
		"address": "Test Address"
	}`

	r := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing required fields: [is_headquarter]")
}
