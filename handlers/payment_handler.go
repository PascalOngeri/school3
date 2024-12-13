package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	
)

// Payment represents a payment structure
type Payment struct {
	PaymentName string
	Term1       float64
	Term2       float64
	Term3       float64
	Total       float64
	Area        string // for bus payments
	Form        string // for compulsory payments
	ID          int
	ClassName   string

	Adm string

	SNo     int
	RegNo   string
	Date    string
	Amount  float64
	Balance float64
	Tot     float64 // Total amount paid
	Balo    float64
}

// PageData represents the data passed to the template
type PageData struct {
	CompulsoryPayments []Payment
	OptionalPayments   []Payment
	BusPayments        []Payment
}

// FetchPayments retrieves payments from the database based on the payment type
func FetchPayments(db *sql.DB, paymentType string) ([]Payment, error) {
	
	var query string
	if paymentType == "compulsory" {
		query = "SELECT paymentname, term1, term2, term3, form, id, (term1 + term2 + term3) AS total FROM feepay"
	} else if paymentType == "optional" {
		query = "SELECT type AS paymentname, t1 AS term1, t2 AS term2, t3 AS term3, id, (t1 + t2 + t3) AS total FROM other"
	} else if paymentType == "bus" {
		query = "SELECT area AS paymentname, t1 AS term1, t2 AS term2, t3 AS term3, id, (t1 + t2 + t3) AS total FROM bus"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		// Scan the data into the Payment struct based on payment type
		if paymentType == "compulsory" {
			err := rows.Scan(&p.PaymentName, &p.Term1, &p.Term2, &p.Term3, &p.Form, &p.ID, &p.Total)
			if err != nil {
				return nil, err
			}
		} else if paymentType == "optional" {
			err := rows.Scan(&p.PaymentName, &p.Term1, &p.Term2, &p.Term3, &p.ID, &p.Total)
			if err != nil {
				return nil, err
			}
		} else if paymentType == "bus" {
			err := rows.Scan(&p.PaymentName, &p.Term1, &p.Term2, &p.Term3, &p.ID, &p.Total)
			if err != nil {
				return nil, err
			}
		}
		payments = append(payments, p)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return payments, nil
}

// UpdatePaymentHandler handles the request to update payments and render the template
func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Fetch data for compulsory payments
	compulsoryPayments, err := FetchPayments(db, "compulsory")
	if err != nil {
		http.Error(w, "Failed to fetch compulsory payments", http.StatusInternalServerError)
		return
	}

	// Fetch data for optional payments
	optionalPayments, err := FetchPayments(db, "optional")
	if err != nil {
		http.Error(w, "Failed to fetch optional payments", http.StatusInternalServerError)
		return
	}

	// Fetch data for bus payments
	busPayments, err := FetchPayments(db, "bus")
	if err != nil {
		http.Error(w, "Failed to fetch bus payments", http.StatusInternalServerError)
		return
	}

	// Create PageData struct to hold all the fetched data
	pageData := PageData{
		CompulsoryPayments: compulsoryPayments,
		OptionalPayments:   optionalPayments,
		BusPayments:        busPayments,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/edelete.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render the template with the data
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
