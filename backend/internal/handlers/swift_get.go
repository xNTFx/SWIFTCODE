package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/validation"
)

// getSwiftCodeDetailsHandler handles GET requests for a single SWIFT code.
func (h *Handler) GetSwiftCodeDetailsHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/v1/swift-codes/"))
	if swiftCode == "" {
		writeJSONError(w, http.StatusBadRequest, "SWIFT code is required")
		return
	}

	if !validation.SwiftCodeRegex.MatchString(swiftCode) {
		writeJSONError(w, http.StatusBadRequest, "Invalid SWIFT code format – it must be exactly 11 letters or digits.")
		return
	}

	swiftCode = strings.ToUpper(swiftCode)
	isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

	if isHeadquarter {
		h.handleHeadquarterSwiftCode(w, swiftCode)
	} else {
		h.handleBranchSwiftCode(w, swiftCode)
	}
}

// supports the SWIFT code for the bank's headquarters.
func (h *Handler) handleHeadquarterSwiftCode(w http.ResponseWriter, swiftCode string) {
	var headquarter models.SwiftCodeHeadquarter
	var branches sql.NullString

	err := h.DB.QueryRow(`
		SELECT 
			sc.address, b.name AS bank_name, c.iso2_code AS country_iso2, 
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
	`, swiftCode).Scan(
		&headquarter.Address, &headquarter.BankName, &headquarter.CountryISO2,
		&headquarter.CountryName, &headquarter.IsHeadquarter, &headquarter.SwiftCode, &branches,
	)

	if err != nil {
		handleDBError(w, err)
		return
	}

	if branches.Valid {
		json.Unmarshal([]byte(branches.String), &headquarter.Branches)
	}

	respondWithJSON(w, http.StatusOK, headquarter)
}

// supports SWIFT code for bank branches
func (h *Handler) handleBranchSwiftCode(w http.ResponseWriter, swiftCode string) {
	var branch models.SwiftCodeBranch

	err := h.DB.QueryRow(`
		SELECT 
			sc.swift_code, b.name AS bank_name, sc.address, 
			c.iso2_code AS country_iso2, c.name AS country_name, sc.is_headquarter
		FROM swift_codes sc
		JOIN banks b ON sc.bank_id = b.id
		JOIN countries c ON b.country_id = c.id
		WHERE sc.swift_code = $1 AND sc.is_headquarter = false;
	`, swiftCode).Scan(
		&branch.SwiftCode, &branch.BankName, &branch.Address,
		&branch.CountryISO2, &branch.CountryName, &branch.IsHeadquarter,
	)

	if err != nil {
		handleDBError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, branch)
}
