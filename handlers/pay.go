package handlers

import (
	"database/sql"
	"fmt"
	"log"
	
	"net/http"
	"strconv"
	"strings"
)

// HandlePayment processes the payment and updates the database
func HandlePayment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the form values for admission number and amount
	adm := r.FormValue("adm")
	amountStr := r.FormValue("ammount")

	// Validate inputs
	if adm == "" || amountStr == "" {
		http.Error(w, "Admission number and amount are required", http.StatusBadRequest)
		return
	}

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Fetch student details from the database
	var fee, balance float64
	var fname, lname, class, phoneNumber string
	stmt, err := db.Prepare("SELECT fee, fname, lname, class, phone FROM registration WHERE adm = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Error preparing database query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(adm).Scan(&fee, &fname, &lname, &class, &phoneNumber)
	if err != nil {
		log.Printf("Error fetching student data: %v", err)
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Calculate the new balance after payment
	balance = fee - amount

	// Setup API call data
	apiURL := "https://infinityschools.xyz/p/api.php"
	data := map[string]interface{}{
		"publicApi": "ISpublic_Api_Keysitq2v5mutip95ra.shabanet",
		"Token":     "ISSecrete_Token_Keya8x3xi4z32959rt1.shabanet",
		"Phone":     phoneNumber,
		"username":  "Pascal Ongeri",
		"password":  "2222",
		"Amount":    amount,
	}

	// Send API request to process payment
	response, err := sendApiRequest(apiURL, data)
	if err != nil || response == "" {
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Process the API response
	if !contains(response, "success: Payment received Successful") {
		http.Error(w, "API call failed", http.StatusInternalServerError)
		return
	}

	// Extract payment reference from the response
	var reference string
	_, err = fmt.Sscanf(response, "from %s amount %f Reference %s", &phoneNumber, &amount, &reference)
	if err != nil {
		http.Error(w, "Error processing payment reference", http.StatusInternalServerError)
		return
	}

	// Insert payment details into the database
	stmt, err = db.Prepare("INSERT INTO payment (adm, amount, bal, reference, paytype) VALUES (?, ?, ?, ?, 'MPESA')")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Error preparing database insert", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(adm, amount, balance, reference)
	if err != nil {
		log.Printf("Error inserting payment record: %v", err)
		http.Error(w, "Error recording payment", http.StatusInternalServerError)
		return
	}

	// Update the student's fee balance in the registration table
	stmt, err = db.Prepare("UPDATE registration SET fee = ? WHERE adm = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Error preparing database update", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(balance, adm)
	if err != nil {
		log.Printf("Error updating student balance: %v", err)
		http.Error(w, "Error updating balance", http.StatusInternalServerError)
		return
	}

	// Log the payment activity
	logMessage := fmt.Sprintf("Recorded payment of school fees for admission number %s. Fee balance is %.2f", adm, balance)
	stmt, err = db.Prepare("INSERT INTO logs (user, activities) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
	}
	_, err = stmt.Exec("admin", logMessage) // Assuming 'admin' is the logged-in user
	if err != nil {
		log.Printf("Error logging payment: %v", err)
	}

	// Respond with success message
	w.Write([]byte("Payment processed successfully"))
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
// Utility function to check if a substring exists in a string
func contains(s, substr string) bool {
	return strings.Index(s, substr) != -1
}

// Function to send the API request (implement it as needed)
func sendApiRequest(url string, data map[string]interface{}) (string, error) {
	// Implement API request logic here, e.g., using net/http
	// Example:
	// resp, err := http.Post(url, "application/json", data)
	// return response as string
	return "", nil
}
