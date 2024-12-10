package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"strings"
	"time"

	"database/sql"
	"strconv"
)
const (
	authURL     = "https://auth.smsafrica.tech/auth/api-key"
	smsURL      = "https://sms-service.smsafrica.tech/message/send/transactional"
	callbackURL = "https://callback.io/123/dlr"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type SmsRequest struct {
	Message     string `json:"message"`
	Msisdn      string `json:"msisdn"`
	SenderID    string `json:"sender_id"`
	CallbackURL string `json:"callback_url"`
}
func SendSms(phone, message string) error {
	// Credentials for API token retrieval

	username := "+254708107995"                                                    // Replace with your username
	password := "e05b0e0d42c608dd08151cfc325da68f1eadd7bf60e457a043bc2e1de39635e2" // Replace with your password
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

	// Step 1: Retrieve API Token
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return fmt.Errorf("error creating token request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64Credentials)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error retrieving token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to get token. Status: %d, Response: %s", resp.StatusCode, body)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("error parsing token response: %v", err)
	}

	apiToken := tokenResponse.Token
	fmt.Println("Retrieved API Token:", apiToken)

	// Step 2: Send SMS
	smsRequest := SmsRequest{
		Message:     message,
		Msisdn:      phone,
		SenderID:    "SMSAFRICA",
		CallbackURL: callbackURL,
	}

	postData, err := json.Marshal(smsRequest)
	if err != nil {
		return fmt.Errorf("error marshalling SMS request: %v", err)
	}

	req, err = http.NewRequest("POST", smsURL, bytes.NewBuffer(postData))
	if err != nil {
		return fmt.Errorf("error creating SMS request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiToken)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending SMS: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading SMS response: %v", err)
	}

	fmt.Printf("SMS Response: %s\n", body)
	return nil
}
// Function to make an API request
func makeRequest(apiurl string, data map[string]string) ([]byte, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(dataJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}

// Payment processing function
func ProcessPayment(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	// Extract userID, totalCost, and phone from the request parameters (URL or body)
	userID := r.FormValue("adm")          // Extract userID from form data
	totalCostStr := r.FormValue("ammount") // Extract totalCost as string
	phone := r.FormValue("phone")         // Extract phone number

	// Convert totalCost to float64
	totalCost, err := strconv.ParseFloat(totalCostStr, 64)
	if err != nil {
		http.Error(w, "Invalid total cost value", http.StatusBadRequest)
		return fmt.Errorf("invalid totalCost value: %v", err)
	}

	apiurl := "https://infinityschools.xyz/p/api.php"

	// Format the total cost as a whole number
	amountStr := fmt.Sprintf("%d", int(totalCost))
	fmt.Println("Processing payment of amount:", amountStr)

	// Data for payment processing
	data := map[string]string{
		"publicApi": "ISpublic_Api_Keysitq2v5mutip95ra.shabanet", // Partner ID
		"Token":     "ISSecrete_Token_Keya8x3xi4z32959rt1.shabanet",
		"Phone":     phone,            // Use the phone number passed as a parameter
		"username":  "Pascal Ongeri",  // Username
		"password":  "2222",           // Password
		"Amount":    amountStr,
	}

	// Log the complete request data
	fmt.Printf("Sending payment request: %+v\n", data)

	// Make the payment request with retries
	var responseBody []byte
	var paymentErr error
	retries := 3
	for i := 0; i < retries; i++ {
		responseBody, paymentErr = makeRequest(apiurl, data)
		if paymentErr == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if paymentErr != nil {
		http.Error(w, fmt.Sprintf("All payment attempts failed: %v", paymentErr), http.StatusInternalServerError)
		return paymentErr
	}

	fmt.Println("Payment API Response:", string(responseBody))

	// Check if the response is a string (non-JSON response)
	responseString := string(responseBody)
	if len(responseString) > 0 && responseString[0] == 's' {
		// The response starts with "success", process it accordingly
		parts := strings.Split(responseString, " ")
		if len(parts) > 5 && parts[0] == "success:" {
			reference := parts[7] 

			var currentFee float64
			err = db.QueryRow("SELECT fee FROM registration WHERE adm = ?", userID).Scan(&currentFee)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching current fee: %v", err), http.StatusInternalServerError)
				return fmt.Errorf("error fetching current fee: %v", err)
			}

			// Calculate the new fee balance
			newFee := currentFee - totalCost
			// Success message: Update the fee in the database
			_, err = db.Exec("UPDATE registration SET fee = fee - ? WHERE adm = ?", totalCost, userID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error updating fee: %v", err), http.StatusInternalServerError)
				return fmt.Errorf("error updating fee: %v", err)
			}

			// Insert payment record into the database
			_, err = db.Exec("INSERT INTO payment (adm, amount, bal) VALUES (?, ?, ?)", userID, totalCost, newFee)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error inserting payment record: %v", err), http.StatusInternalServerError)
				return fmt.Errorf("error inserting payment record: %v", err)
			}

			smsMessage := fmt.Sprintf("Payment received successfully. Amount: %.2f. New balance: %.2f. Reference: %s", totalCost, newFee, reference)
			smsErr := SendSms(phone, smsMessage)
			if smsErr != nil {
				http.Error(w, fmt.Sprintf("Payment successful but failed to send SMS: %v", smsErr), http.StatusInternalServerError)
				return fmt.Errorf("error sending SMS: %v", smsErr)
			}
			// Send success response with details
			w.WriteHeader(http.StatusOK)
			//fmt.Fprintf(w, "Payment received successfully from %s, amount %.2f. Reference: %s", phone, totalCost, reference)
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
			
		
		} else {
			http.Error(w, "Error: Payment not received", http.StatusBadRequest)
		}
	} else if len(responseString) > 0 && responseString[0] == 'e' {
		// If the response starts with "error", check for cancellation
		if responseString == "error: Payment Cancelled by user" {
			// Only update status to 2 if the payment is canceled
			

			// Send failure response with cancellation message
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "Payment canceled by user. Status updated to 2.")
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
		} else {
			// Handle other error messages
			http.Error(w, fmt.Sprintf("Payment failed: %v", responseString), http.StatusBadRequest)
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
	
		}
	} else {
		// If the response is not in expected string format, it could be a JSON response
		var apiResponse map[string]interface{}
		if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing API response: %v", err), http.StatusInternalServerError)
			return fmt.Errorf("error parsing API response: %v", err)
		}

		// Check if the payment was successful based on JSON response
		if status, ok := apiResponse["status"].(string); !ok || status != "success" {
			http.Error(w, fmt.Sprintf("Payment failed: %v", apiResponse), http.StatusBadRequest)
			return fmt.Errorf("payment failed: %v", apiResponse)
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
		}

		// Send response for successful JSON response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Payment processed successfully for user %s", userID)
	}

	return nil
}
