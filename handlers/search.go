package handlers

import (
	"html/template"
	"log"
	"net/http"

)

type STU struct {
	Adm      string
	Fname    string
	Mname    string
	Lname    string
	Gender   string
	Faname   string
	Maname   string
	Class    string
	Phone    string
	Phone1   string
	Address  string
	Email    string
	Fee      string
	T1       string
	T2       string
	T3       string
	Dob      string
	Image    string
	Username string
	Password string
}

func add1(i int) int {
	return i + 1
}
func searchStudentHandler(w http.ResponseWriter, r *http.Request) {
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
	funcMap := template.FuncMap{
		"add1": add1, // Register the add1 function
	}
	// Handle POST request (form submission)
	if r.Method == http.MethodPost {
		searchData := r.FormValue("searchdata")
		if searchData == "" {
			http.Error(w, "Please enter a search term", http.StatusBadRequest)
			return
		}

		// Query the database to search for students by their admission number (Adm)
		rows, err := db.Query("SELECT adm, fname, mname, lname, gender, faname, maname, class, phone, phone1, address, email, fee, t1, t2, t3, dob, image, username, password FROM registration WHERE adm LIKE ?", "%"+searchData+"%")
		if err != nil {
			http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create a slice to hold the students found
		var students []STU
		for rows.Next() {
			var student STU
			if err := rows.Scan(&student.Adm, &student.Fname, &student.Mname, &student.Lname, &student.Gender, &student.Faname, &student.Maname, &student.Class, &student.Phone, &student.Phone1, &student.Address, &student.Email, &student.Fee, &student.T1, &student.T2, &student.T3, &student.Dob, &student.Image, &student.Username, &student.Password); err != nil {
				log.Println(err)
				continue
			}
			students = append(students, student)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading from the database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Register the custom template function

		// Render the result template with the students data
		tmpl, err := template.New("search").Funcs(funcMap).ParseFiles(
			"templates/search.html", // Update this path as needed
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v", err)
			return
		}

		// Pass the students data to the template
		// Pass the students data to the template
		err = tmpl.Execute(w, students) // 'students' is the slice of STU passed as context
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}

	} else {
		// Handle GET request (render search form)
		tmpl, err := template.ParseFiles(
			"templates/search.html", // Update this path as needed
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v", err)
			return
		}

		// Execute the template (empty data for initial search page)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	}
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}