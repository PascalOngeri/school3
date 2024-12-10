
package handlers
import (
	"database/sql"
	"html/template"
	"net/http"
	
	"log"
	"fmt"
	
	
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
// HandleLogin handles login requests
func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var foundInAdmin bool
		var adm, phone, role string
		var fee float64

		// Authenticate user in tbladmin
		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			foundInAdmin = true
			role = "admin"
		} else {
			// Authenticate user in registration table
			queryRegistration := "SELECT id, adm, username, phone, password, fee FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password, &fee)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			role = "user"
			log.Printf("User ID: %d, Adm: %s, Username: %s, Phone: %s, Fee: %f", userID, adm, username, phone, fee)
		}

		// Redirect based on role
		if foundInAdmin {
			redirectURL := "/dashboard?role=admin&userID=" + fmt.Sprintf("%d", userID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else if role == "user" {

			// User (student) should go to the parent section
			redirectURL := "/parent?role=user&userID=" + fmt.Sprintf("%d", userID) +
				"&adm=" + adm + "&username=" + username + "&phone=" + phone + "&fee=" + fmt.Sprintf("%f", fee)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		return
	}

	// Render the login page for GET requests
	api, _ := getAPIDetails(db)
	renderLoginPage(w, api, "")
}
