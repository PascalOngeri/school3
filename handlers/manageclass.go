package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	
)

func Manageclass(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Fetch data from the database
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
	users := []User{}
	rows, err := db.Query("SELECT id, class, t1, t2, t3, fee FROM classes")
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error during db.Query: %v\n", err) // Debug log
		return
	}
	defer rows.Close()

	// Process rows
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Class, &user.T1, &user.T2, &user.T3, &user.Fee); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error during rows.Scan: %v\n", err) // Debug log
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error during rows iteration: %v\n", err) // Debug log
		return
	}

	// Parse the template files
	tmpl, err := template.ParseFiles(
		"templates/manage-class.html",
		"includes/footer.html",
		"includes/header.html",
		"includes/sidebar.html",
	)
	if err != nil {
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v\n", err) // Debug log
		return
	}

	// Execute the template with the fetched data
	err = tmpl.Execute(w, users)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err) // Debug log
		return
	}
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}