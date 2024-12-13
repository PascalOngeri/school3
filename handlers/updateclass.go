package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	
)

func UpdateClass(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			// Pata data kutoka kwenye fomu
			id := r.FormValue("id")
			className := r.FormValue("className")
			t1Fee := r.FormValue("t1Fee")
			t2Fee := r.FormValue("t2Fee")
			t3Fee := r.FormValue("t3Fee")

			// Badilisha fees kuwa float
			t1, err := strconv.ParseFloat(t1Fee, 64)
			if err != nil {
				log.Printf("Failed to parse term 1 fee: %v", err)
				http.Error(w, "Invalid term 1 fee", http.StatusBadRequest)
				return
			}

			t2, err := strconv.ParseFloat(t2Fee, 64)
			if err != nil {
				log.Printf("Failed to parse term 2 fee: %v", err)
				http.Error(w, "Invalid term 2 fee", http.StatusBadRequest)
				return
			}

			t3, err := strconv.ParseFloat(t3Fee, 64)
			if err != nil {
				log.Printf("Failed to parse term 3 fee: %v", err)
				http.Error(w, "Invalid term 3 fee", http.StatusBadRequest)
				return
			}

			// Hesabu totalFee
			totalFee := t1 + t2 + t3

			// Sasisha darasa kwenye database
			_, err = db.Exec("UPDATE classes SET class = ?, t1 = ?, t2 = ?, t3 = ?, fee = ? WHERE id = ?",
				className, t1, t2, t3, totalFee, id)
			if err != nil {
				log.Printf("Failed to update class: %v", err)
				http.Error(w, "Failed to update class.", http.StatusInternalServerError)
				return
			}

			// Rejea kwa ukurasa wa kusimamia darasa
			http.Redirect(w, r, "/manage", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}