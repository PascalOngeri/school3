package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"net/http"
	"strings"
)

// SelectPhonesHandler fetches phone numbers and returns them as a comma-separated string
func SelectPhonesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Fetch all phone numbers from the `registration` table
		rows, err := db.Query("SELECT phone FROM registration")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch phone numbers: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var phones []string
		for rows.Next() {
			var phone string
			if err := rows.Scan(&phone); err != nil {
				http.Error(w, fmt.Sprintf("Failed to scan phone number: %v", err), http.StatusInternalServerError)
				return
			}
			phones = append(phones, phone)
		}

		// Check for iteration errors
		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to iterate over rows: %v", err), http.StatusInternalServerError)
			return
		}

		// Join the phone numbers with commas
		phoneNumbers := ""
		if len(phones) > 0 {
			phoneNumbers = strings.Join(phones, ", ")
		}

		// Return the phone numbers as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"phoneNumbers": phoneNumbers})
	}
}