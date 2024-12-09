package handlers

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Student struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Class      string

	Gender           string
	DOB              string
	AdmissionNumber  string
	Image            string
	FatherName       string
	MotherName       string
	ContactNumber    string
	AltContactNumber string
	Address          string
	UserName         string
	Password         string
}

type Class struct {
	ID   int
	Name string

	Class string
	T1    float64
	T2    float64
	T3    float64
	Fee   float64
}

func GetClassDetails(db *sql.DB, class string) (float64, float64, float64, float64, error) {
	// Query to fetch details
	query := `SELECT t1, t2, t3, fee FROM classes WHERE class = ?`
	var t1, t2, t3, fee float64

	// Execute the query
	err := db.QueryRow(query, class).Scan(&t1, &t2, &t3, &fee)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, 0, nil // No data found for the given class
		}
		return 0, 0, 0, 0, err // Handle other errors
	}

	return t1, t2, t3, fee, nil
}
func Addstudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles(
		"templates/addstudent.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch classes from the database
	var classes []Class
	rows, err := db.Query("SELECT id, class FROM classes")
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Collect form data
		student := Student{
			FirstName:        r.FormValue("fname"),
			MiddleName:       r.FormValue("mname"),
			LastName:         r.FormValue("lname"),
			Email:            r.FormValue("stuemail"),
			Class:            r.FormValue("stuclass"),
			Gender:           r.FormValue("gender"),
			DOB:              r.FormValue("dob"),
			AdmissionNumber:  r.FormValue("stuid"),
			FatherName:       r.FormValue("faname"),
			MotherName:       r.FormValue("maname"),
			ContactNumber:    r.FormValue("connum"),
			AltContactNumber: r.FormValue("altconnum"),
			Address:          r.FormValue("address"),
			UserName:         r.FormValue("uname"),
			Password:         hashPassword(r.FormValue("password")),
		}

		// Handle file upload
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			uploadDir := "uploads/"
			os.MkdirAll(uploadDir, os.ModePerm)

			// Generate a random file name to avoid collisions
			fileName := filepath.Join(uploadDir, randomFileName(handler.Filename))
			destFile, err := os.Create(fileName)
			if err == nil {
				defer destFile.Close()
				_, _ = destFile.ReadFrom(file)
				student.Image = fileName
			}
		}
		t1, t2, t3, fee, err := GetClassDetails(db, student.Class)
		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Failed to fetch class details: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert data into the database
		query := `
			INSERT INTO registration (
				adm, fname, mname, lname, gender, faname, maname, class, phone, phone1,
				address, email, fee, t1, t2, t3, dob, image, username, password
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = db.Exec(query,
			student.AdmissionNumber, student.FirstName, student.MiddleName, student.LastName,
			student.Gender, student.FatherName, student.MotherName, student.Class,
			student.ContactNumber, student.AltContactNumber, student.Address,
			student.Email, fee, t1, t2, t3, student.DOB, student.Image,
			student.UserName, student.Password,
		)
		if err != nil {
			http.Error(w, "Error inserting student: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Student %s successfully added to the database", student.FirstName)
		http.Redirect(w, r, "/addstudent", http.StatusSeeOther)
		return
	}

	// Render the form
	data := map[string]interface{}{
		"Title":   "Add Students",
		"Classes": classes,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to hash a password
func hashPassword(password string) string {
	// Add your hashing logic here (e.g., bcrypt)
	return password // Replace with hashed password
}

// Helper function to generate a random file name
func randomFileName(original string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%x_%s", buf, original)
}
