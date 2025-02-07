package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func handleDBError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Resource not found")
	} else {
		writeJSONError(w, http.StatusInternalServerError, "Database query failed")
		log.Printf("Database error: %v", err)
	}
}
