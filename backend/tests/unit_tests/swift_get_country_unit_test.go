package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"backend/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestGetSwiftCodesByCountryHandler_InvalidFormat verifies that an invalid country ISO2 code format returns a bad request error.
func TestGetSwiftCodesByCountryHandler_InvalidFormat(t *testing.T) {
	t.Log("Testing invalid country ISO2 code format returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/INVALID!@#", nil)
	w := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid Country ISO2 Code format – it must be exactly 2 letters."}`, w.Body.String())
}

// TestGetSwiftCodesByCountryHandler_NotFound verifies that requesting a non-existent country returns a 404 error.
func TestGetSwiftCodesByCountryHandler_NotFound(t *testing.T) {
	t.Log("Testing retrieval of SWIFT codes for a non-existent country returns 404 Not Found")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT c.name AS country_name,
		COALESCE(json_agg(json_build_object(
			'bankName', b.name,
			'address', sc.address,
			'countryISO2', c.iso2_code,
			'isHeadquarter', sc.is_headquarter,
			'swiftCode', sc.swift_code
		)), '[]'::json) AS swift_codes
		FROM countries c
		LEFT JOIN banks b ON b.country_id = c.id
		LEFT JOIN swift_codes sc ON sc.bank_id = b.id
		WHERE c.iso2_code = $1
		GROUP BY c.name;
	`)).WithArgs("XX").WillReturnError(sql.ErrNoRows)

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/XX", nil)
	w := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestGetSwiftCodesByCountryHandler_Success verifies that requesting a valid country returns a list of SWIFT codes.
func TestGetSwiftCodesByCountryHandler_Success(t *testing.T) {
	t.Log("Testing retrieval of SWIFT codes for a valid country")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT c.name AS country_name,
		COALESCE(json_agg(json_build_object(
			'bankName', b.name,
			'address', sc.address,
			'countryISO2', c.iso2_code,
			'isHeadquarter', sc.is_headquarter,
			'swiftCode', sc.swift_code
		)), '[]'::json) AS swift_codes
		FROM countries c
		LEFT JOIN banks b ON b.country_id = c.id
		LEFT JOIN swift_codes sc ON sc.bank_id = b.id
		WHERE c.iso2_code = $1
		GROUP BY c.name;
	`)).WithArgs("PL").WillReturnRows(sqlmock.NewRows([]string{"country_name", "swift_codes"}).AddRow("Poland", `[
		{"bankName": "Test Bank", "address": "Test Address", "countryISO2": "PL", "isHeadquarter": true, "swiftCode": "ABCDEFGHXXX"}
	]`))

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/PL", nil)
	w := httptest.NewRecorder()

	handler.GetSwiftCodesByCountryHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"countryISO2":"PL","countryName":"Poland","swiftCodes":[{"bankName":"Test Bank","address":"Test Address","countryISO2":"PL","isHeadquarter":true,"swiftCode":"ABCDEFGHXXX"}]}`, w.Body.String())
}
