package handlers

import (

  
    "net/http"
    "html/template"
    "database/sql"
   "log"

    // Import the function to send SMS
    // Update this import path to where your SMS sending functions are defined
)

func Send(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
    // Handle POST request
    if r.Method == http.MethodPost {
        phone := r.FormValue("phone")
        message := r.FormValue("message")

        if phone == "" || message == "" {
            http.Error(w, "Phone number and message are required", http.StatusBadRequest)
            return
        }

        

        // Send the SMS
        SendSms(phone, message)

        // After sending the SMS, return success message
        http.Redirect(w, r, "/send", http.StatusSeeOther) // Redirect after form submission
        return
    }

    // Parse the template files
    tmpl, err := template.ParseFiles("templates/send.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Data to pass to the template
    data := map[string]interface{}{
        "Title": "Send SMS", // Example dynamic data
    }

    // Execute the template and write to the response
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}else {
        // If role is not recognized, redirect to login
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}