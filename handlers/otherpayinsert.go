package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// FormHandler handles form submission and database operations
func Insert(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	if r.Method == http.MethodPost {
		// Retrieve form values
		payname := r.FormValue("payname")
		term1Str := r.FormValue("term1")
		term2Str := r.FormValue("term2")
		term3Str := r.FormValue("term3")

		// Parse term values to float
		term1, err := strconv.ParseFloat(term1Str, 64)
		if err != nil {
			log.Println("Error parsing term1:", err)
			http.Error(w, "Invalid input for Term 1", http.StatusBadRequest)
			return
		}

		term2, err := strconv.ParseFloat(term2Str, 64)
		if err != nil {
			log.Println("Error parsing term2:", err)
			http.Error(w, "Invalid input for Term 2", http.StatusBadRequest)
			return
		}

		term3, err := strconv.ParseFloat(term3Str, 64)
		if err != nil {
			log.Println("Error parsing term3:", err)
			http.Error(w, "Invalid input for Term 3", http.StatusBadRequest)
			return
		}

		// Calculate total amount
		amount := term1 + term2 + term3

		// Prepare SQL query
		query := `
			INSERT INTO other (type, t1, t2, t3, amount)
			VALUES (?, ?, ?, ?, ?)`

		// Execute the SQL query
		_, err = db.Exec(query, payname, term1, term2, term3, amount)
		if err != nil {
			log.Println("Database insertion error:", err)
			http.Error(w, "Failed to save data. Please try again later.", http.StatusInternalServerError)
			return
		}

		// Redirect to confirmation page or another page
		http.Redirect(w, r, "/setfee", http.StatusSeeOther)
		return
	}

	// If not POST, show the form again
	http.Redirect(w, r, "/setfee", http.StatusSeeOther)
}