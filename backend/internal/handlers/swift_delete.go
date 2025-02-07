package handlers

import (
	"net/http"
	"strings"

	"backend/internal/validation"
)

// deleteSwiftCodeHandler obsługuje żądania DELETE usuwające SWIFT code.
func (h *Handler) DeleteSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
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

	query := `SELECT delete_swift_code($1)`
	// CREATE OR REPLACE FUNCTION delete_swift_code(swift_code_input VARCHAR(11))
	// RETURNS BOOLEAN AS $$
	// DECLARE
	// 	bank_id_var UUID;
	// 	country_id_var UUID;
	// BEGIN
	// 	DELETE FROM swift_codes
	// 	WHERE swift_code = swift_code_input
	// 	RE

	// 	IF bank_id_var IS NULL THE
	
	// 	IF bank_id_var IS NULL THEN
	// 		RETURN FALSE;
	// 	END IF;
	
	// 	IF NOT EXISTS (SELECT 1 FROM swift_codes WHERE bank_id = bank_id_var) THEN
	// 		SELECT country_id INTO country_id_var FROM banks WHERE id = bank_id_var;
	// 		DELETE FROM banks WHERE id = bank_id_var;
	
	// 		IF NOT EXISTS (SELECT 1 FROM banks WHERE country_id = country_id_var) THEN
	// 			DELETE FROM countries WHERE id = country_id_var;
	// 		END IF;
	// 	END IF;
	
	// 	RETURN TRUE;
	// END;
	// $$ LANGUAGE plpgsql;	

	var swiftDeleted bool
	err := h.DB.QueryRow(query, swiftCode).Scan(&swiftDeleted)
	if err != nil {
		handleDBError(w, err)
		return
	}

	if !swiftDeleted {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"message": "SWIFT code not found, nothing to delete"})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "SWIFT code deleted successfully"})
}
