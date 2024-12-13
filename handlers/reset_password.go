package handlers

import (
	"database/sql"
	"html/template"
	"log"
	
	"net/http"
)

// ResetPasswordHandler handles the password reset process
func ResetPasswordHandler(db *sql.DB) http.HandlerFunc {
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
		if r.Method == http.MethodGet {
			// Render the password reset form
			tmpl, err := template.ParseFiles("templates/reset.html")
			if err != nil {
				log.Printf("Error loading templates: %v", err)
				http.Error(w, "Error loading templates", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			email := r.FormValue("email")
			mobile := r.FormValue("mobile")
			newPassword := r.FormValue("newpassword")
			confirmPassword := r.FormValue("confirmpassword")

			// Validate passwords
			if newPassword != confirmPassword {
				http.Error(w, "Passwords do not match", http.StatusBadRequest)
				return
			}

			// Check if email and mobile exist in the database
			var existingEmail string
			query := `SELECT Email FROM tbladmin WHERE Email = ? AND MobileNumber = ?`
			err := db.QueryRow(query, email, mobile).Scan(&existingEmail)
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or mobile number", http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("Error checking user: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			// Update the password in the database
			updateQuery := `UPDATE tbladmin SET Password = ? WHERE Email = ? AND MobileNumber = ?`
			_, err = db.Exec(updateQuery, newPassword, email, mobile)
			if err != nil {
				log.Printf("Error updating password: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			// Redirect to a success page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}