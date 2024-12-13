package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Class structure

// EditClass handler for editing class details
func EditClass(db *sql.DB) http.HandlerFunc {
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
		// Get the edit ID from the URL
		editID := r.URL.Query().Get("editid")
		if editID == "" {
			http.Error(w, "Missing edit ID parameter", http.StatusBadRequest)
			return
		}

		// Define the class variable
		var class Class

		// Fetch class details from the database
		err := db.QueryRow("SELECT id, class, t1, t2, t3, fee FROM classes WHERE id = ?", editID).
			Scan(&class.ID, &class.Class, &class.T1, &class.T2, &class.T3, &class.Fee)

		if err != nil {
			log.Printf("Failed to fetch class: %v", err)
			http.Error(w, "Failed to fetch class details.", http.StatusInternalServerError)
			return
		}

		// Parse the template
		tmpl, err := template.ParseFiles("templates/edit-class.html", "includes/header.html", "includes/sidebar.html", "includes/footer.html")
		if err != nil {
			log.Printf("Template parsing failed: %v", err)
			http.Error(w, "Failed to load page templates.", http.StatusInternalServerError)
			return
		}

		// Execute the template with class data
		data := map[string]interface{}{
			"Title": "Edit Class",
			"Class": class,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
		}
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}