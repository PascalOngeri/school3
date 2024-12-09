package handlers

import (
	"database/sql"
	"net/http"
)

// Function ya kufuta mwanafunzi
func DeleteStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM registration WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}
}
