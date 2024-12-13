package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteNotice deletes a public notice by its ID
func DeleteNotice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		// Get the ID from the query parameter
		delID := r.URL.Query().Get("delID")
		if delID == "" {
			http.Error(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		// Execute the DELETE query
		_, err := db.Exec("DELETE FROM tblpublicnotice WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to delete notice: %v", err)
			http.Error(w, "Failed to delete notice.", http.StatusInternalServerError)
			return
		}

		// Redirect to the manage public notice page after deletion
		http.Redirect(w, r, "/manage-public-notice", http.StatusSeeOther)
	} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}
