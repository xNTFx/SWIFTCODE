package tests

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handlers"

	"github.com/stretchr/testify/assert"
)

// TestDeleteSwiftCodeHandler_Success verifies that an existing SWIFT code is deleted successfully.
func TestDeleteSwiftCodeHandler_Success(t *testing.T) {
	t.Log("Testing successful SWIFT code deletion")

	// Ensure SWIFT code exists before deletion
	_, err := db.Exec(`
		INSERT INTO countries (iso2_code, name) VALUES ('PL', 'Poland')
		ON CONFLICT (iso2_code) DO NOTHING;
		INSERT INTO banks (name, country_id) VALUES ('Test Bank', (SELECT id FROM countries WHERE iso2_code = 'PL'))
		ON CONFLICT (name, country_id) DO NOTHING;
		INSERT INTO swift_codes (swift_code, bank_id, is_headquarter, address)
		VALUES ('ABCDEFGHXXX', (SELECT id FROM banks WHERE name = 'Test Bank' AND country_id = (SELECT id FROM countries WHERE iso2_code = 'PL')), true, 'Test Address')
		ON CONFLICT (swift_code) DO NOTHING;
	`)
	assert.NoError(t, err)

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/ABCDEFGHXXX", nil)
	rec := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code deleted successfully", response["message"])

	// Verify that SWIFT code is deleted
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM swift_codes WHERE swift_code = 'ABCDEFGHXXX')`).Scan(&exists)
	assert.NoError(t, err)
	assert.False(t, exists, "SWIFT code should be deleted from database")
}

// TestDeleteSwiftCodeHandler_NotFound verifies that trying to delete a non-existent SWIFT code returns 404.
func TestDeleteSwiftCodeHandler_NotFound(t *testing.T) {
	t.Log("Testing deletion of a non-existent SWIFT code returns 404 Not Found")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/NONEXISTENT", nil)
	rec := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code not found, nothing to delete", response["message"])
}

// TestDeleteSwiftCodeHandler_InvalidFormat verifies that an invalid SWIFT code format results in a 400 error.
func TestDeleteSwiftCodeHandler_InvalidFormat(t *testing.T) {
	t.Log("Testing invalid SWIFT code format returns 400 Bad Request")

	handler := handlers.NewHandler(db)
	req := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/INVALID_CODE", nil)
	rec := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid SWIFT code format – it must be exactly 11 letters or digits.")
}

// TestDeleteSwiftCodeHandler_DatabaseError verifies that a database error results in a 500 response.
func TestDeleteSwiftCodeHandler_DatabaseError(t *testing.T) {
	t.Log("Testing database error returns 500 Internal Server Error")

	// Simulate database failure by using a new instance with an invalid DSN
	brokenDB, _ := sql.Open("postgres", "postgres://invalid:invalid@localhost:5432/invalid?sslmode=disable")
	handler := handlers.NewHandler(brokenDB)

	req := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/ABCDEFGHXXX", nil)
	rec := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Database query failed")

	// Ensure that the original database connection is still valid
	err := db.Ping()
	assert.NoError(t, err)
}
