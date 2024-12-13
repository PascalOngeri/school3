package handlers

import (
	"html/template"
	
	"net/http"
	"log"
)

func report(w http.ResponseWriter, r *http.Request) {
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
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/report.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		// Handle the error properly, e.g., by returning a 500 status
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title": "Manage Class", // Example dynamic data
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		// Handle the error properly
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}