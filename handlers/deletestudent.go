package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// Function ya kufuta mwanafunzi
func DeleteStudent(db *sql.DB) http.HandlerFunc {

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
		// Get the student ID from the query parameter
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		// Execute the DELETE query
		_, err := db.Exec("DELETE FROM registration WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting student: %v", err)
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		// Redirect to the manage student page after deletion
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}
}
