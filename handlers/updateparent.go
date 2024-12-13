package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// User represents the user structure.


// SettingsHandler handles user settings (GET to display and POST to update).
func SettingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Admission Number from the form or URL.
		adm := r.FormValue("adm")
		if adm == "" {
			http.Error(w, "Admission number is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Fetch user data from the database.
			var user User
			err := db.QueryRow("SELECT adm, username, password, phone FROM registration WHERE adm = ?", adm).Scan(&user.AdmissionNumber, &user.Username, &user.Password, &user.Phone)
			if err != nil {
				log.Printf("Error fetching user data: %v", err)
				http.Error(w, "Error fetching user data", http.StatusInternalServerError)
				return
			}

			// Render the template with user data.
			tmpl, err := template.ParseFiles("templates/parent.html")
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, user)

		case http.MethodPost:
			// Retrieve form values for update.
			username := r.FormValue("username")
			password := r.FormValue("password")
			phone := r.FormValue("phone")

			if username == "" || password == "" || phone == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			// Update the database with new values.
			_, err := db.Exec("UPDATE registration SET username = ?, password = ?, phone = ? WHERE adm = ?", username, password, phone, adm)
			if err != nil {
				log.Printf("Error updating user data: %v", err)
				http.Error(w, "Error updating user data", http.StatusInternalServerError)
				return
			}

			// Redirect to the settings page after update.
			http.Redirect(w, r, "/login", http.StatusSeeOther)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}