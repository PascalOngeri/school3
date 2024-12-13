package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

)

// Notice represents a public notice with a title and message
type Notice struct {
	Title   string
	Message string
	Date    string
	ID      interface{}
}

// ManagePubNot handles public notices by fetching them from the database
func ManagePubNot(db *sql.DB) http.HandlerFunc {
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
		// Query to fetch public notices
		query := "SELECT ID, NoticeTitle, NoticeMessage FROM tblpublicnotice"
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Database query failed: %v", err)
			http.Error(w, "Failed to fetch notices from the database.", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notices []Notice
		// Process rows
		for rows.Next() {
			var notice Notice
			if err := rows.Scan(&notice.ID, &notice.Title, &notice.Message); err != nil {
				log.Printf("Error scanning row: %v", err)
				http.Error(w, "Error processing data from the database.", http.StatusInternalServerError)
				return
			}
			notices = append(notices, notice)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
			http.Error(w, "Error reading notices from the database.", http.StatusInternalServerError)
			return
		}

		// Parse template files
		tmpl, err := template.ParseFiles(
			"templates/managepubnot.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("Template parsing failed: %v", err)
			http.Error(w, "Failed to load page templates.", http.StatusInternalServerError)
			return
		}

		// Pass data to the template
		data := map[string]interface{}{
			"Title":   "Manage Public Notice",
			"Notices": notices,
		}

		// Render the template
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
		}
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}