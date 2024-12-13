package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteClass deletes a class by ID
func DeleteClass(db *sql.DB) http.HandlerFunc {
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
		// Retrieve the `delid` query parameter
		delID := r.URL.Query().Get("delid")
		if delID == "" {
			http.Error(w, "Missing class ID to delete.", http.StatusBadRequest)
			return
		}

		// Delete the class from the database
		_, err := db.Exec("DELETE FROM classes WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to delete class with ID %s: %v", delID, err)
			http.Error(w, "Failed to delete class.", http.StatusInternalServerError)
			return
		}

		// Redirect back to the manage page
		http.Redirect(w, r, "/manage", http.StatusSeeOther)
	} 
	

}
}
