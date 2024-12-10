package handlers

import (
	"net/http"

)

// LogoutHandler handles the logout process for JWT-based authentication
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the JWT by setting an expired cookie
	

		// Optionally, you can add any other steps needed before logout (like logging out from sessions, etc.)

		// Redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
