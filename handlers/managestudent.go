package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Struct kwa data ya mwanafunzi
type SelectStudent struct {
	ID    int
	Adm   string
	Class string
	Fname string
	Mname string
	Lname string
	Fee   float64
	Email string
	Phone string
}

// Function ya kusimamia wanafunzi
func ManageStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var sele []SelectStudent

		rows, err := db.Query("SELECT id, adm, class, fname, mname, lname, fee, email, phone FROM registration")
		if err != nil {
			http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error during db.Query: %v\n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var student SelectStudent
			if err := rows.Scan(&student.ID, &student.Adm, &student.Class, &student.Fname, &student.Mname, &student.Lname, &student.Fee, &student.Email, &student.Phone); err != nil {
				http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error during rows.Scan: %v\n", err)
				return
			}
			sele = append(sele, student)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error during rows iteration: %v\n", err)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/managestudent.html",
			"includes/footer.html",
			"includes/header.html",
			"includes/sidebar.html",
		)
		if err != nil {
			http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v\n", err)
			return
		}

		err = tmpl.Execute(w, sele)
		if err != nil {
			http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v\n", err)
			return
		}
	}
}
