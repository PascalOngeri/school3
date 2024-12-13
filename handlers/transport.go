package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	
)

// TemplateData holds data to be passed to the template
type TemplateData struct {
	Message            string
	Total              float64
	CompulsoryPayments []Payment
	OptionalPayments   []Payment
	BusPayments        []Payment
	Title              string
	Username           string
	AdmissionNumber    string
	Password           string
	Phone              string
	Payments           []Payment
	Notices            []Notice
}

// FormHandler handles form submission and database operations
func FormHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	role := roleCookie.Value
	//userID := r.URL.Query().Get("userID")
	// If role is "admin", show the dashboard
	if role == "admin" {
	if r.Method == http.MethodPost {
		// Retrieve form values
		term1Str := r.FormValue("term1")
		term2Str := r.FormValue("term2")
		term3Str := r.FormValue("term3")
		area := r.FormValue("area")

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

		// Calculate total
		total := term1 + term2 + term3

		// Insert data into the database
		query := `
			INSERT INTO bus (area, t1, t2, t3, amount)
			VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(query, area, term1, term2, term3, total)
		if err != nil {
			log.Println("Database insertion error:", err)
			http.Error(w, "Failed to save data. Please try again later.", http.StatusInternalServerError)
			return
		}

		// Redirect to confirmation or the original form
		http.Redirect(w, r, "/setfee", http.StatusSeeOther)
		return
	}

	// If not POST, redirect to the form
	http.Redirect(w, r, "/setfee", http.StatusSeeOther)
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}