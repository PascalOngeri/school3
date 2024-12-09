package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Class    string
	T1       string
	T2       string
	T3       string
	Fee      string
	Adm      string
	UserName string
	Phone    string
	Password string
	Address  string
	Phone2   string
	Phone1   string
	MotherN  string
	FatherN  string
	Image    string
	Dob      string
	Gender   string
	Email    string
	Lname    string
	Mname    string
	Fname    string

	Date       string
	User       string
	Activities string
}

// Function to get user by ID
func getUserByEmail(db *sql.DB, id string) (User, error) {
	var user User
	err := db.QueryRow("SELECT adm, fname, mname, lname, gender, faname, maname, class, phone, phone1, address, email, fee, t1, t2, t3, dob, image, username, password FROM registration WHERE adm = ?", id).Scan(&user.Adm, &user.Fname, &user.Mname, &user.Lname, &user.Gender, &user.FatherN, &user.MotherN,
		&user.Class, &user.Phone1, &user.Phone2, &user.Address, &user.Email, &user.Fee, &user.T1,
		&user.T2, &user.T3, &user.Dob, &user.Image, &user.UserName, &user.Password)
	return user, err
}

// Handler to update user details
func UpdateUserFormHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Handle GET request to display the form
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		// Fetch user details based on the ID
		user, err := getUserByEmail(db, id)
		if err != nil {
			log.Printf("Error fetching user details: %v", err)
			http.Error(w, "Error fetching user details", http.StatusInternalServerError)
			return
		}

		// Parse and render the template with user data
		tmpl, err := template.ParseFiles(
			"templates/updatestudent.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}

		// Render template with user data
		if err := tmpl.Execute(w, user); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle POST request to update the user data
	if r.Method == "POST" {
		// Get form values
		email := r.FormValue("stuemail")
		username := r.FormValue("uname")
		password := r.FormValue("password")
		fname := r.FormValue("fname")
		mname := r.FormValue("mname")
		lname := r.FormValue("lname")
		class := r.FormValue("class")
		gender := r.FormValue("gender")
		dob := r.FormValue("dob")
		adm := r.FormValue("stuid")
		faname := r.FormValue("faname")
		maname := r.FormValue("maname")
		connum := r.FormValue("connum")
		altconnum := r.FormValue("altconnum")
		address := r.FormValue("address")

		// Validate required fields
		if email == "" || username == "" || fname == "" || lname == "" || class == "" {
			http.Error(w, "All required fields must be filled", http.StatusBadRequest)
			return
		}

		// Hash password if provided
		var hashedPassword []byte
		if password != "" {
			var err error
			hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Error hashing password: %v", err)
				http.Error(w, "Error processing password", http.StatusInternalServerError)
				return
			}
		}

		// Use the hashed password or the existing password
		updatePassword := hashedPassword
		if password == "" {
			// Use the current password if no new one is provided
			updatePassword = nil
		}

		// Update user details in the database
		_, err := db.Exec(
			"UPDATE registration SET fname=?, mname=?, lname=?, gender=?, faname=?, maname=?, class=?, phone=?, phone1=?, address=?, email=?, username=?, password=?, dob=? WHERE adm=?",
			fname, mname, lname, gender, faname, maname, class, connum, altconnum, address, email, username, updatePassword, dob, adm,
		)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Error updating user details", http.StatusInternalServerError)
			return
		}

		// Redirect after successful update
		http.Redirect(w, r, "/managestudent?success=1", http.StatusSeeOther)
	}
}
