package handlers

import (
	"net/http"
	"time"

)


// LogoutHandler handles the logout process for JWT-based authentication
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the JWT by setting an expired cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "role",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "userID",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "adm",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "phone",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "fee",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to a past date to delete the cookie
		HttpOnly: true,
	})

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
