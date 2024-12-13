package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	
)

// UpdateBusPaymentHandler handles the update of bus payment details
func UpdateBusPaymentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if r.Method == "GET" {
		// Get the ID of the bus payment to edit from the URL
		id := r.URL.Query().Get("id")

		// Fetch current bus payment details from the database
		var area, t1, t2, t3 string
		var amount float64
		query := `SELECT area, t1, t2, t3, amount FROM bus WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&area, &t1, &t2, &t3, &amount)
		if err != nil {
			http.Error(w, "Error fetching bus payment details", http.StatusInternalServerError)
			log.Printf("Error fetching bus payment: %v", err)
			return
		}

		// Prepare data for template
		data := struct {
			ID     string
			Area   string
			T1     string
			T2     string
			T3     string
			Amount float64
		}{
			ID:     id,
			Area:   area,
			T1:     t1,
			T2:     t2,
			T3:     t3,
			Amount: amount,
		}

		// Render the template with current bus payment data
		tmpl, err := template.ParseFiles("templates/editB.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)

	} else if r.Method == "POST" {
		// Handle form submission

		id := r.FormValue("id")
		area := r.FormValue("fname")
		t1 := r.FormValue("lname")
		t2 := r.FormValue("stuemail")
		t3 := r.FormValue("dob")

		// Convert the form values for amount if necessary (assuming it's a number)
		amount := r.FormValue("amount")

		// Update the bus payment in the database
		updateQuery := `UPDATE bus SET area = ?, t1 = ?, t2 = ?, t3 = ?, amount = ? WHERE id = ?`
		_, err := db.Exec(updateQuery, area, t1, t2, t3, amount, id)
		if err != nil {
			http.Error(w, "Error updating bus payment", http.StatusInternalServerError)
			log.Printf("Error updating bus payment: %v", err)
			return
		}

		// Redirect to the success page or another page
		http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
	}
} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}