package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions" // Ensure this is installed via go get
)

// Global session store
var store = sessions.NewCookieStore([]byte("your-secret-key"))

// Dashboard handles the /dashboard route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session
	// session, _ := store.Get(r, "store")
	// if session.Values["sturecmsaid"] == nil {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }

	// Parse templates
	tmpl, err := template.ParseFiles(
		"templates/dashboard.html",
		"includes/footer.html",
		"includes/header.html",
		"includes/sidebar.html",
	)
	if err != nil {
		// Handle the error properly
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
