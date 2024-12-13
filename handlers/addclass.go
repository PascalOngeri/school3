package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	


)

// Replace this with your actual secret key


// ValidateJWT function for validating the token

func AddClass(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Validate JWT and get the claims
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
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get the class name from the form
		className := r.FormValue("cname")
		if className == "" {
			http.Error(w, "Class name is required", http.StatusBadRequest)
			return
		}

		// Insert the class into the database
		_, err := db.Exec("INSERT INTO classes (class, fee, t1, t2, t3) VALUES (?,?,?,?,?)", className, 0, 0, 0, 0)
		if err != nil {
			http.Error(w, "Failed to add class: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error inserting class into database: %v", err)
			return
		}

		// Redirect to a confirmation page or reload the form with a success message
		http.Redirect(w, r, "/addclass", http.StatusSeeOther)
		return
	}

	// Render the template
	tmpl, err := template.ParseFiles(
		"templates/addclass.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}else if role == "user" {
		// If the role is "user", redirect to the parent section
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}