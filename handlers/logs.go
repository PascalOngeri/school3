package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	
)

func Logs(db *sql.DB) http.HandlerFunc {
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
		rows, err := db.Query("SELECT id, date, user, activities FROM logs")
		if err != nil {
			http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error during db.Query: %v\n", err)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Date, &user.User, &user.Activities); err != nil {
				http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error during rows.Scan: %v\n", err)
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error during rows iteration: %v\n", err)
			return
		}

		// Pass data to the template
		tmpl, err := template.ParseFiles(
			"templates/logs.html",
			"includes/footer.html",
			"includes/header.html",
			"includes/sidebar.html",
		)
		if err != nil {
			http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v\n", err)
			return
		}

		data := map[string]interface{}{
			"Users": users,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v\n", err)
			return
		}
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}