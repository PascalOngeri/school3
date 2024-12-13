package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	
	"os"
	"path/filepath"
)

// SettingsHandler handles settings updates
func SettingsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		handlePostRequest(w, r, db)
		return
	}

	// Handle GET request
	tmpl, err := template.ParseFiles(
		"templates/setting.html", // Use a relevant name
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Render the settings page
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// handlePostRequest handles the POST logic
func handlePostRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the form
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		log.Printf("Form parsing error: %v", err)
		return
	}

	// Handle file upload
	filePath, err := saveUploadedFile(r)
	if err != nil {
		http.Error(w, "File upload failed", http.StatusInternalServerError)
		log.Printf("File upload error: %v", err)
		return
	}

	// Get school name
	schoolName := r.FormValue("name")
	if schoolName == "" {
		http.Error(w, "School name is required", http.StatusBadRequest)
		log.Println("School name is missing")
		return
	}

	// Update the database
	query := "UPDATE api SET icon = ?, name = ?"
	_, err = db.Exec(query, filePath, schoolName)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}

	log.Printf("Updated settings: School Name - %s, Logo Path - %s", schoolName, filePath)

	// Redirect to settings page
	http.Redirect(w, r, "/setting", http.StatusSeeOther)
}

// saveUploadedFile handles the file upload and saves it to the server
func saveUploadedFile(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Validate file type (optional)
	if !validateFileType(handler) {
		return "", http.ErrNotSupported
	}

	// Ensure upload directory exists
	uploadDir := "assets/images/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Save file
	filePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	return filePath, err
}

// validateFileType ensures the uploaded file is an image
func validateFileType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range allowedTypes {
		if fileHeader.Header.Get("Content-Type") == t {
			return true
		}
	}
	return false
}
