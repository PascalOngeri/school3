package handlers

import (
	"database/sql"
	"log"
	"net/http"
	
)

// DeleteCompulsoryHandler handles the deletion of a record from the "feepay" table
func DeleteCompulsoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get the payment name from the query parameters
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
	paymentName := r.URL.Query().Get("delid")
	if paymentName == "" {
		http.Error(w, "Missing PaymentName", http.StatusBadRequest)
		return
	}

	// Prepare the delete statement
	query := "DELETE FROM feepay WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the delete statement
	result, err := stmt.Exec(paymentName)
	if err != nil {
		log.Println("Error executing delete:", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	if rowsAffected == 0 {
		http.Error(w, "No record found to delete", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
}

}