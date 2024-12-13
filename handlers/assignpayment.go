package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	
)

// PaymentOption represents a payment option
type PaymentOption struct {
	ID   int
	Name string
}

// AreaOption represents an area option
type AreaOption struct {
	ID   int
	Name string
}

// AssignPaymentsPageData holds data for rendering the assign payments page
type AssignPaymentsPageData struct {
	Payments  []PaymentOption
	Areas     []AreaOption
	Durations []string // Added durations for terms
}

// HandleAssignPayments handles rendering the assign payments page
func HandleAssignPayments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Fetch payment options
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
	paymentRows, err := db.Query("SELECT id, name FROM payments")
	if err != nil {
		log.Printf("Error fetching payments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer paymentRows.Close()

	var payments []PaymentOption
	for paymentRows.Next() {
		var payment PaymentOption
		if err := paymentRows.Scan(&payment.ID, &payment.Name); err != nil {
			log.Printf("Error scanning payment: %v", err)
			continue
		}
		payments = append(payments, payment)
	}

	// Fetch area options
	areaRows, err := db.Query("SELECT id, name FROM areas")
	if err != nil {
		log.Printf("Error fetching areas: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer areaRows.Close()

	var areas []AreaOption
	for areaRows.Next() {
		var area AreaOption
		if err := areaRows.Scan(&area.ID, &area.Name); err != nil {
			log.Printf("Error scanning area: %v", err)
			continue
		}
		areas = append(areas, area)
	}

	// Define payment durations
	durations := []string{"Term 1", "Term 2", "Term 3", "Term 2 and Term 3", "All"}

	// Prepare data for the template
	data := AssignPaymentsPageData{
		Payments:  payments,
		Areas:     areas,
		Durations: durations, // Include durations
	}

	// Render template
	tmpl, err := template.ParseFiles(
		"templates/optionalpay.html",
		"templates/header.html",
		"templates/sidebar.html",
	)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
} 
}