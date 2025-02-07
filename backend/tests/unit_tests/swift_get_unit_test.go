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

// TestGetSwiftCodeDetailsHandler_InvalidFormat verifies that an invalid SWIFT code format returns a bad request error.
func TestGetSwiftCodeDetailsHandler_InvalidFormat(t *testing.T) {
	t.Log("Testing invalid SWIFT code format returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/INVALID!@#", nil)
	w := httptest.NewRecorder()

	handler.SwiftHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid SWIFT code format – it must be exactly 11 letters or digits."}`, w.Body.String())
}

// TestGetSwiftCodeDetailsHandler_HeadquarterFound verifies that retrieving a headquarter SWIFT code works correctly.
func TestGetSwiftCodeDetailsHandler_HeadquarterFound(t *testing.T) {
	t.Log("Testing retrieval of a headquarter SWIFT code")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT sc.address, b.name AS bank_name, c.iso2_code AS country_iso2,
		c.name AS country_name, sc.is_headquarter, sc.swift_code,
		(
			SELECT COALESCE(json_agg(json_build_object(
				'bankName', b2.name,
				'address', sw.address,
				'countryISO2', c2.iso2_code,
				'isHeadquarter', sw.is_headquarter,
				'swiftCode', sw.swift_code
			)), '[]'::json)
			FROM swift_codes sw
			JOIN banks b2 ON sw.bank_id = b2.id
			JOIN countries c2 ON b2.country_id = c2.id
			WHERE LEFT(sw.swift_code, 8) = LEFT($1, 8)
			AND sw.swift_code != $1
			AND sw.is_headquarter = false
		) AS branches
		FROM swift_codes sc
		JOIN banks b ON sc.bank_id = b.id
		JOIN countries c ON b.country_id = c.id
		WHERE sc.swift_code = $1;
	`)).WithArgs("ABCDEFGHXXX").WillReturnRows(sqlmock.NewRows([]string{
		"address", "bank_name", "country_iso2", "country_name", "is_headquarter", "swift_code", "branches",
	}).AddRow("Test Address", "Test Bank", "PL", "Poland", true, "ABCDEFGHXXX", "[]"))

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/ABCDEFGHXXX", nil)
	w := httptest.NewRecorder()

	handler.SwiftHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetSwiftCodeDetailsHandler_BranchFound verifies that retrieving a branch SWIFT code works correctly.
func TestGetSwiftCodeDetailsHandler_BranchFound(t *testing.T) {
	t.Log("Testing retrieval of a branch SWIFT code")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT sc.swift_code, b.name AS bank_name, sc.address,
		c.iso2_code AS country_iso2, c.name AS country_name, sc.is_headquarter
		FROM swift_codes sc
		JOIN banks b ON sc.bank_id = b.id
		JOIN countries c ON b.country_id = c.id
		WHERE sc.swift_code = $1 AND sc.is_headquarter = false;
	`)).WithArgs("ABCDEFGH001").WillReturnRows(sqlmock.NewRows([]string{
		"swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter",
	}).AddRow("ABCDEFGH001", "Test Bank", "Branch Address", "PL", "Poland", false))

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/ABCDEFGH001", nil)
	w := httptest.NewRecorder()

	handler.SwiftHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetSwiftCodeDetailsHandler_NotFound verifies that requesting a non-existent SWIFT code returns a 404 error.
func TestGetSwiftCodeDetailsHandler_NotFound(t *testing.T) {
	t.Log("Testing retrieval of a non-existent SWIFT code returns 404 Not Found")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT sc.swift_code, b.name AS bank_name, sc.address,
		c.iso2_code AS country_iso2, c.name AS country_name, sc.is_headquarter
		FROM swift_codes sc
		JOIN banks b ON sc.bank_id = b.id
		JOIN countries c ON b.country_id = c.id
		WHERE sc.swift_code = $1 AND sc.is_headquarter = false;
	`)).WithArgs("NONEXISTENT").WillReturnError(sql.ErrNoRows)

	r := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/NONEXISTENT", nil)
	w := httptest.NewRecorder()

	handler.SwiftHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"Resource not found"}`, w.Body.String())
}
