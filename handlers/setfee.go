package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	
)

// Class represents a class structure for dropdown data

// SetFeeHandler handles the Set Fee page
func SetFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse templates

	tmpl, err := template.ParseFiles(
		"templates/regfee.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Unable to load templates", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	// Variables to store form values and messages
	var message string
	var total float64

	// Handle form submission
	if r.Method == http.MethodPost {
		// Retrieve form values
		class := r.FormValue("class")
		payName := r.FormValue("payname")
		term1 := r.FormValue("term1")
		term2 := r.FormValue("term2")
		term3 := r.FormValue("term3")

		// Convert term values to float
		term1Value, err1 := strconv.ParseFloat(term1, 64)
		term2Value, err2 := strconv.ParseFloat(term2, 64)
		term3Value, err3 := strconv.ParseFloat(term3, 64)

		// Check for conversion errors
		if err1 != nil || err2 != nil || err3 != nil {
			message = "Invalid term values provided."
		} else {
			// Calculate total
			total = term1Value + term2Value + term3Value

			// Insert or update fee in the database
			err := InsertOrUpdateFee(db, class, payName, term1Value, term2Value, term3Value, total)
			if err == nil {
				message = "Fee processed successfully!"
			} else {
				message = "Error processing fee: " + err.Error()
			}
		}
	}

	// Fetch classes for dropdown
	var classes []Class
	rows, err := db.Query("SELECT id, class FROM classes")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		classes = append(classes, class)
	}
	log.Println("Classes fetched:", classes)

	// Prepare data for the template
	data := map[string]interface{}{
		"Title":   "Set Fee",
		"Classes": classes,
		"Message": message,
		"Total":   total,
	}

	// Render the template
	if err := tmpl.ExecuteTemplate(w, "regfee.html", data); err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

// InsertOrUpdateFee inserts or updates a fee record in the database
func InsertOrUpdateFee(db *sql.DB, class, payName string, term1, term2, term3, total float64) error {
	// Insert query for feepay table
	insertQuery := `
        INSERT INTO feepay (form, paymentname, term1, term2, term3, amount)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(insertQuery, class, payName, term1, term2, term3, total)
	updateQuery := `
            UPDATE classes
            SET 
                t1 = t1 + ?, 
                t2 = t2 + ?, 
                t3 = t3 + ?, 
                fee = fee + ?
            WHERE class = ?
        `
	_, err = db.Exec(updateQuery, term1, term2, term3, total, class)
	if err != nil {
		log.Println("Insert into feepay failed, attempting update:", err)

		// Update query for classes table

		if err != nil {
			log.Println("Update classes failed:", err)
			return err
		}
		log.Println("Classes updated successfully")
	}
	return nil
}
