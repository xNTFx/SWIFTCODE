package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/validation"
)

// getswiftcodesbycountryhandler handles GET requests for a single country ISO2 code.
func (h *Handler) GetSwiftCodesByCountryHandler(w http.ResponseWriter, r *http.Request) {
	countryISO2Code := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/v1/swift-codes/country/"))
	if countryISO2Code == "" {
		writeJSONError(w, http.StatusBadRequest, "Country ISO2 Code is required")
		return
	}

	if !validation.CountryIsoRegex.MatchString(countryISO2Code) {
		writeJSONError(w, http.StatusBadRequest, "Invalid Country ISO2 Code format – it must be exactly 2 letters.")
		return
	}

	countryISO2Code = strings.ToUpper(countryISO2Code)

	var countrySwiftCodes models.SwiftCodeByCountryISO2
	var swiftCodes sql.NullString
	var countryName sql.NullString

	err := h.DB.QueryRow(`
		SELECT 
			c.name AS country_name,
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
	`, countryISO2Code).Scan(&countryName, &swiftCodes)

	if err != nil {
		handleDBError(w, err)
		return
	}

	countrySwiftCodes.CountryISO2 = countryISO2Code
	countrySwiftCodes.CountryName = countryName.String

	if swiftCodes.Valid {
		json.Unmarshal([]byte(swiftCodes.String), &countrySwiftCodes.SwiftCodes)
	} else {
		countrySwiftCodes.SwiftCodes = []models.SwiftCodeDetails{}
	}

	respondWithJSON(w, http.StatusOK, countrySwiftCodes)
}
