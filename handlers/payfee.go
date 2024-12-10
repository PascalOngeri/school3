package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// PayFeeHandler handles fee payment logic
func PayFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Retrieve session


	if r.Method == http.MethodPost {
		adm := r.FormValue("adm")
		amount := r.FormValue("ammount")

		// Validate admission number
		if adm == "" {
			http.Error(w, "Admission number is required", http.StatusBadRequest)
			return
		}

		// Validate and convert amount
		if amount == "" {
			http.Error(w, "Amount is required", http.StatusBadRequest)
			return
		}
		amt, err := strconv.ParseFloat(amount, 64)
		if err != nil || amt <= 0 {
			log.Printf("Invalid amount format: %s", amount)
			http.Error(w, "Invalid amount format. Please enter a positive number, e.g., 2000.00", http.StatusBadRequest)
			return
		}

		// Update student fee
		sqlUpdate := "UPDATE registration SET fee = fee - ? WHERE adm = ?"
		result, err := db.Exec(sqlUpdate, amt, adm)
		if err != nil {
			log.Printf("Error updating fee: %v", err)
			http.Error(w, "Error updating fee. Please try again later.", http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			log.Printf("No student found with adm: %s", adm)
			http.Error(w, "Admission number not found.", http.StatusNotFound)
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

	// Fetch recent payments
	rows, err := db.Query("SELECT id, adm, date, amount, bal, (SELECT SUM(amount) FROM payment) AS total_amount, (SELECT SUM(bal) FROM payment) AS total_balance FROM payment ORDER BY id DESC")
	if err != nil {
		log.Println("Error fetching payments:", err)
		http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.ID, &p.Adm, &p.Date, &p.Amount, &p.Balance,&p.Tot,&p.Balo)
		if err != nil {
			log.Println("Error scanning payment:", err)
			continue
		}
		payments = append(payments, p)
	}
	classRows, err := db.Query("SELECT id, class FROM classes ")
	if err != nil {
		log.Println("Error fetching classes:", err)
		http.Error(w, "Failed to fetch classes", http.StatusInternalServerError)
		return
	}
	defer classRows.Close()

	var classes []struct {
		ID   int
		Name string
	}
	for classRows.Next() {
		var cls struct {
			ID   int
			Name string
		}
		err := classRows.Scan(&cls.ID, &cls.Name)
		if err != nil {
			log.Println("Error scanning class:", err)
			continue
		}
		classes = append(classes, cls)
	}

	// Prepare data for template rendering
	data := struct {
		Payments []Payment
		Classes  []struct {
			ID   int
			Name string
		}
	}{
		Payments: payments,
		Classes:  classes,
	}

	// Load and render templates
	tmpl, err := template.ParseFiles(
		"templates/payfee.html",

		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
