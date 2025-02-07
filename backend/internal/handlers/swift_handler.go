package handlers

import (
	"database/sql"
	"net/http"
)

// Handler is a structure that stores a reference to the database.
type Handler struct {
	DB *sql.DB
}

// NewHandler creates a new handler with a reference to the database.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// SwiftHandler handles HTTP requests to the /v1/swift-codes/ endpoint.
func (h *Handler) SwiftHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSwiftCodeDetailsHandler(w, r)
	case http.MethodPost:
		h.PostSwiftCodeHandler(w, r)
	case http.MethodDelete:
		h.DeleteSwiftCodeHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
