package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteClass deletes a class by ID
func DeleteClass(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delID := r.URL.Query().Get("delid")
		_, err := db.Exec("DELETE FROM classes WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to delete class: %v", err)
			http.Error(w, "Failed to delete class.", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/manage", http.StatusSeeOther)
	}
}
