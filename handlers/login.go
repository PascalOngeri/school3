package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// Set the session cookie name and expiration
	sessionCookieName = "user_id"
	sessionDuration   = 30 * time.Minute
)

// API structure
type API struct {
	Name  string
	Icon  string
	IName string
}

// LoginData structure for rendering the login page
type LoginData struct {
	Name     string
	Icon     string
	Username string
	Password string
}

// Get API details from the database
func getAPIDetails(db *sql.DB) (API, error) {
	var api API
	query := "SELECT name, icon, iname FROM api LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&api.Name, &api.Icon, &api.IName)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		return api, err
	}
	return api, nil
}

// Render the login page
func renderLoginPage(w http.ResponseWriter, api API, username string) {
	loginData := LoginData{
		Name:     api.Name,
		Icon:     api.Icon,
		Username: username,
		Password: "",
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, loginData)
}

// HandleLogin handles login requests
func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var foundInAdmin bool
		var adm, phone, role,form string
		var fee float64

		// Authenticate user in tbladmin
		queryAdmin := "SELECT ID, UserName, Password FROM tbladmin WHERE UserName = ?"
		var storedPassword string
		err := db.QueryRow(queryAdmin, username).Scan(&userID, &username, &storedPassword)
		if err == nil && password == storedPassword { // Plain text password comparison
			foundInAdmin = true
			role = "admin"
		} else {
			// Authenticate user in registration table
			queryRegistration := "SELECT id, adm, username, phone, password, fee,class  FROM registration WHERE username = ?"
			err = db.QueryRow(queryRegistration, username).Scan(&userID, &adm, &username, &phone, &storedPassword, &fee,&form)
			if err != nil || password != storedPassword { // Plain text password comparison
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			role = "user"
			log.Printf("User ID: %d, Adm: %s, Username: %s, Phone: %s, Fee: %f", userID, adm, username, phone, fee)
		}

		// Set cookies for user details
		http.SetCookie(w, &http.Cookie{
			Name:     "role",
			Value:    role,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "form",
			Value:    form,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "userID",
			Value:    fmt.Sprintf("%d", userID),
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "adm",
			Value:    adm,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    username,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "phone",
			Value:    phone,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "fee",
			Value:    fmt.Sprintf("%f", fee),
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
		})

http.SetCookie(w, &http.Cookie{
    Name:     "Password",
    Value:    storedPassword, // Set the value directly
    Expires:  time.Now().Add(sessionDuration),
    HttpOnly: true,
})

		// Redirect based on role
		if foundInAdmin {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		} else if role == "user" {
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
			return
		}
		return
	}

	// Render the login page for GET requests
	api, _ := getAPIDetails(db)
	renderLoginPage(w, api, "")
}
