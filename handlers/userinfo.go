package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)


// UserInfoHandler handles the GET and POST requests for user info.
func UserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve cookies for role and userID
		roleCookie, err := r.Cookie("role")
		idCookie, err := r.Cookie("userID")
		if err != nil {
			log.Printf("Error retrieving cookies: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		role := roleCookie.Value
		ID := idCookie.Value

		if role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Fetch user info from the database.
			var user User
			err := db.QueryRow("SELECT MobileNumber, Email, UserName FROM tbladmin WHERE ID = ?", ID).Scan(&user.Phone, &user.Email, &user.Username)
			if err != nil {
				log.Printf("Error fetching user info: %v", err)
				http.Error(w, "Error fetching user info", http.StatusInternalServerError)
				return
			}

			// Render the template with user data.
			tmpl, err := template.ParseFiles("templates/userinfo.html")
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, user)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}

		case http.MethodPost:
			// Parse form values
			err := r.ParseForm()
			if err != nil {
				log.Printf("Error parsing form: %v", err)
				http.Error(w, "Invalid form submission", http.StatusBadRequest)
				return
			}

			phone := r.FormValue("mobile")
			email := r.FormValue("email")
			username := r.FormValue("newpassword")

			// Update user info in the database.
			_, err = db.Exec("UPDATE tbladmin SET MobileNumber = ?, Email = ?, UserName = ? WHERE ID = ?", phone, email, username, ID)
			if err != nil {
				log.Printf("Error updating user info: %v", err)
				http.Error(w, "Error updating user info", http.StatusInternalServerError)
				return
			}

			// Redirect to the GET method after the update.
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
