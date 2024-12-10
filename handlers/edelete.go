package handlers

import (
	"html/template"

	"net/http"
	
)

func edelete(w http.ResponseWriter, r *http.Request) {

	// Parse the template files
	tmpl, err := template.ParseFiles("templates/edelete.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
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
}
