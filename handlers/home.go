package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// HomePageData is the structure used to hold data for rendering the home page template
type HomePageData struct {
	Title           string
	Username        string
	AdmissionNumber string
	Password        string
	Phone           string
	Payments        []Payment
	Notices         []Notice
}

// HomeHandler handles the request for the user home page
func HomeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get cookies for user info
	roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userIDCookie, err := r.Cookie("userID")
	if err != nil {
		log.Printf("Error getting userID cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	admCookie, err := r.Cookie("adm")
	if err != nil {
		log.Printf("Error getting adm cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		log.Printf("Error getting username cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	phoneCookie, err := r.Cookie("phone")
	if err != nil {
		log.Printf("Error getting phone cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	feeCookie, err := r.Cookie("fee")
	if err != nil {
		log.Printf("Error getting fee cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	pass, err := r.Cookie("Password")
if err != nil {
    log.Printf("Error getting password cookie: %v", err)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
}

	// Log or use the cookie values
	role := roleCookie.Value
	userID := userIDCookie.Value
	adm := admCookie.Value
	username := usernameCookie.Value
	phone := phoneCookie.Value
	fee := feeCookie.Value
	passi := pass.Value

	log.Printf("Role: %s, User ID: %s, Adm: %s, Username: %s, Phone: %s, Fee: %s", role, userID, adm, username, phone, fee)
log.Printf("Password cookie retrieved: %s", pass.Value)

	// Check if the role is "user"
	if role == "user" {
		// Fetch payment history
		payments, err := getPayments(db, adm)
		if err != nil {
			log.Printf("Failed to fetch payments: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		// Fetch notices
		notices, err := getNotices(db)
		if err != nil {
			log.Printf("Failed to fetch notices: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		// Prepare data for the template
		data := HomePageData{
			Title:           "Infinityschools Analytics",
			Username:        username,
			AdmissionNumber: adm,
			Password:        passi, // Password is not used here
			Phone:           phone,
			Payments:        payments,
			Notices:         notices,
		}

		// Render the template
		tmpl, err := template.ParseFiles("templates/parent.html", "includes/footer.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	} else {
		// Redirect to login if the role is not "user"
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// getPayments fetches payment history from the database
func getPayments(db *sql.DB, adm string) ([]Payment, error) {
	rows, err := db.Query("SELECT id, adm, date, amount, bal FROM payment WHERE adm = ? ORDER BY id DESC", adm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.SNo, &p.RegNo, &p.Date, &p.Amount, &p.Balance)
		if err != nil {
			log.Printf("Failed to scan payment: %v", err)
			continue
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// getNotices fetches notices from the database
func getNotices(db *sql.DB) ([]Notice, error) {
	rows, err := db.Query("SELECT NoticeTitle, NoticeMessage, CreationDate FROM tblpublicnotice")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []Notice
	for rows.Next() {
		var n Notice
		err := rows.Scan(&n.Title, &n.Message, &n.Date)
		if err != nil {
			log.Printf("Failed to scan notice: %v", err)
			continue
		}
		notices = append(notices, n)
	}
	return notices, nil
}
