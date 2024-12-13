package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	
)

// GetClassDetails retrieves t1, t2, t3, and fee for a specific class from the classes table

// ManageUser handles adding and deleting users
func ManageUser(db *sql.DB) http.HandlerFunc {
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
		// Check if user is logged in
	
		if r.Method == http.MethodPost {

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
				log.Printf("Form parsing error: %v", err)
				return
			}

			action := r.FormValue("submit") // Capture which button was clicked
			if action == "Add" {
				// Add user logic
				AName := r.FormValue("adminname")
				mobno := r.FormValue("mobilenumber")
				email := r.FormValue("email")
				pass := r.FormValue("password")
				username := r.FormValue("username")

				// Validate input
				if AName == "" || mobno == "" || email == "" || pass == "" || username == "" {
					http.Error(w, "All fields are required.", http.StatusBadRequest)
					log.Println("Validation error: missing required fields")
					return
				}

				// Hash the password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
					log.Printf("Password hashing error: %v", err)
					return
				}
 log.Printf("Authenticated user: %s",  hashedPassword)
				// Insert data into the database
				query := `INSERT INTO tbladmin (AdminName, Email, UserName, Password, MobileNumber) VALUES (?, ?, ?, ?, ?)`
				_, err = db.Exec(query, AName, email, username, pass, mobno)
				if err != nil {
					log.Printf("Database insertion error: %v", err)
					http.Error(w, "Failed to add user: "+err.Error(), http.StatusInternalServerError)
					return
				}

				log.Println("User successfully added")
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}

			if action == "Delete" {
				// Delete user logic
				username := r.FormValue("username")

				if username == "" {
					http.Error(w, "Username is required for deletion.", http.StatusBadRequest)
					log.Println("Validation error: username is missing")
					return
				}

				query := `DELETE FROM tblAdmin WHERE UserName = ?`
				result, err := db.Exec(query, username)
				if err != nil {
					log.Printf("Database deletion error: %v", err)
					http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
					return
				}

				rowsAffected, _ := result.RowsAffected()
				if rowsAffected == 0 {
					http.Error(w, "No user found with the provided username.", http.StatusNotFound)
					log.Println("Deletion error: no matching user")
					return
				}

				log.Printf("User %s successfully deleted", username)
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}
		}

		// Render the form template for GET requests
		tmpl, err := template.ParseFiles(
			"templates/adduser.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Template parsing error: %v", err)
			return
		}

		// Render the template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
			log.Printf("Template execution error: %v", err)
		}
	} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}
