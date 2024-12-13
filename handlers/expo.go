package handlers

import (
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"strings"
)

// ExportHandler handles the CSV export request
func ExportHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
    sqlQuery := `SELECT adm, CONCAT(fname, ' ', mname, ' ', lname) AS student_name, class, fee, phone, gender, email, address, dob, faname, maname, username FROM registration`
    rows, err := db.Query(sqlQuery)
    if err != nil {
        log.Println("Error querying the database:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment; filename=student_fee_data.csv")
    w.Header().Set("Cache-Control", "no-store")

    writer := csv.NewWriter(w)
    defer writer.Flush()

    headers := []string{"Admission No.", "Student Name", "Class/Grade/Form", "Fee Balance", "Phone", "Gender", "Email", "Address", "Date of Birth", "Father's Name", "Mother's Name", "Username"}
    if err := writer.Write(headers); err != nil {
        log.Println("Error writing CSV headers:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    for rows.Next() {
        var adm, studentName, className, fee, phone, gender, email, address, dob, fatherName, motherName, username string
        if err := rows.Scan(&adm, &studentName, &className, &fee, &phone, &gender, &email, &address, &dob, &fatherName, &motherName, &username); err != nil {
            log.Println("Error scanning row:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        record := []string{
            adm,
            strings.ToUpper(studentName),
            className,
            fee,
            phone,
            gender,
            email,
            address,
            dob,
            fatherName,
            motherName,
            username,
        }
        if err := writer.Write(record); err != nil {
            log.Println("Error writing CSV record:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }
}}
