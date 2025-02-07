package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestDeleteSwiftCodeHandler_Success verifies that an existing SWIFT code is deleted successfully.
func TestDeleteSwiftCodeHandler_Success(t *testing.T) {
	t.Log("Testing successful SWIFT code deletion")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(`SELECT delete_swift_code\(\$1\)`).
		WithArgs("ABCDEFGHXXX").
		WillReturnRows(sqlmock.NewRows([]string{"delete_swift_code"}).AddRow(true))

	r := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/ABCDEFGHXXX", nil)
	w := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"SWIFT code deleted successfully"}`, w.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSwiftCodeHandler_NotFound verifies that trying to delete a non-existent SWIFT code returns 404.
func TestDeleteSwiftCodeHandler_NotFound(t *testing.T) {
	t.Log("Testing deletion of a non-existent SWIFT code returns 404 Not Found")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(`SELECT delete_swift_code\(\$1\)`).
		WithArgs("NONEXISTENT").
		WillReturnRows(sqlmock.NewRows([]string{"delete_swift_code"}).AddRow(false))

	r := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/NONEXISTENT", nil)
	w := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"SWIFT code not found, nothing to delete"}`, w.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSwiftCodeHandler_InvalidFormat verifies that an invalid SWIFT code format results in a 400 error.
func TestDeleteSwiftCodeHandler_InvalidFormat(t *testing.T) {
	t.Log("Testing invalid SWIFT code format returns 400 Bad Request")
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	r := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/INVALID_CODE", nil)
	w := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid SWIFT code format – it must be exactly 11 letters or digits.")
}

// TestDeleteSwiftCodeHandler_DatabaseError verifies that a database error results in a 500 response.
func TestDeleteSwiftCodeHandler_DatabaseError(t *testing.T) {
	t.Log("Testing database error returns 500 Internal Server Error")
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHandler(db)

	mock.ExpectQuery(`SELECT delete_swift_code\(\$1\)`).
		WithArgs("ABCDEFGHXXX").
		WillReturnError(sql.ErrConnDone)

	r := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/ABCDEFGHXXX", nil)
	w := httptest.NewRecorder()

	handler.DeleteSwiftCodeHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Database query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}
