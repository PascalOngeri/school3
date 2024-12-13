package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

)

func AddPubNot(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve form values
		nottitle := r.FormValue("nottitle")
		notmsg := r.FormValue("notmsg")

		// Log the received form data
		log.Printf("Notice Title: %s, Notice Message: %s", nottitle, notmsg)

		// Check if form data is valid
		if nottitle == "" || notmsg == "" {
			http.Error(w, "Notice Title and Message are required fields.", http.StatusBadRequest)
			return
		}

		// Insert data into the database
		_, err := db.Exec("INSERT INTO tblpublicnotice (NoticeTitle, NoticeMessage) VALUES (?, ?)", nottitle, notmsg)
		if err != nil {
			log.Printf("Failed to insert notice: %v", err) // Log the error
			http.Error(w, "Failed to insert notice: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Notice successfully added")

		// Redirect to the form page (or any other success page)
		http.Redirect(w, r, "/addpubnot", http.StatusSeeOther)
		return
	}

	// Render the template for GET requests
	tmpl, err := template.ParseFiles(
		"templates/addpubnotice.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}