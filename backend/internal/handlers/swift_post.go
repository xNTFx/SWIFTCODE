package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/validation"
)

// postswiftcodehandler handles post requests adding new swift code.
func (h *Handler) PostSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	var body models.SwiftCodeBranch

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	// removal of whitespace characters
	body.SwiftCode = strings.TrimSpace(body.SwiftCode)
	body.BankName = strings.TrimSpace(body.BankName)
	body.CountryISO2 = strings.TrimSpace(body.CountryISO2)
	body.CountryName = strings.TrimSpace(body.CountryName)
	body.Address = strings.TrimSpace(body.Address)

	// check of required fields
	missingFields := validation.ValidateSwiftCodeFields(body)
	if len(missingFields) > 0 {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Missing required fields: %v", missingFields))
		return
	}

	// input validation
	validationErrors := validation.ValidateSwiftCodeBranch(body)
	if len(validationErrors) > 0 {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation errors: %v", validationErrors))
		return
	}

	// conversion to uppercase
	body.SwiftCode = strings.ToUpper(body.SwiftCode)
	body.BankName = strings.ToUpper(body.BankName)
	body.CountryISO2 = strings.ToUpper(body.CountryISO2)
	body.CountryName = strings.ToUpper(body.CountryName)
	body.Address = strings.ToUpper(body.Address)

	isHeadquarter := strings.HasSuffix(body.SwiftCode, "XXX")
	if isHeadquarter != *body.IsHeadquarter {
		writeJSONError(w, http.StatusBadRequest, "Mismatch between SWIFT code format and headquarter status")
		return
	}

	query := `
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
	`

	var insertedCount int
	err := h.DB.QueryRow(query,
		body.CountryISO2,
		body.CountryName,
		body.BankName,
		body.SwiftCode,
		*body.IsHeadquarter,
		body.Address,
	).Scan(&insertedCount)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert SWIFT code")
		log.Printf("Insert error: %v", err)
		return
	}

	if insertedCount == 0 {
		writeJSONError(w, http.StatusConflict, "SWIFT code already exists")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "SWIFT code added successfully"})
}
