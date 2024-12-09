package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func PayFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		adm := r.FormValue("adm")
		amount := r.FormValue("ammount")

		// Ukaguzi wa thamani tupu
		if amount == "" {
			http.Error(w, "Amount is required", http.StatusBadRequest)
			return
		}

		// Jaribu kubadilisha `amount` kuwa `float64`
		amt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Printf("Invalid amount format: %s", amount)
			http.Error(w, "Invalid amount format. Please enter a valid number, e.g., 2000.00", http.StatusBadRequest)
			return
		}

		// Update student fee
		sqlUpdate := "UPDATE registration SET fee = fee - ? WHERE adm = ?"
		_, err = db.Exec(sqlUpdate, amt, adm)
		if err != nil {
			log.Printf("Error updating fee: %v", err)
			http.Error(w, "Error updating fee. Please try again later.", http.StatusInternalServerError)
			return
		}

		// Insert payment record
		sqlInsert := "INSERT INTO payment (adm, amount, bal) VALUES (?, ?, ?)"
		v := 0.0
		_, err = db.Exec(sqlInsert, adm, amt, v)
		if err != nil {
			log.Printf("Error inserting payment: %v", err)
			http.Error(w, "Error recording payment. Please try again later.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/payfee?success=true", http.StatusSeeOther)
		return
	}

	// Render template for GET request
	tmpl, err := template.ParseFiles(
		"templates/payfee.html",
		"includes/footer.html",
		"includes/header.html",
		"includes/sidebar.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
