
package handlers

import (

	
	"fmt"

	"log"
)



func encryptURLParameters(role string, userID int, adm string, username string, phone string, fee float64) string {
	// Define a key for AES encryption (32 bytes for AES-256)
	key := []byte("a very secret key12345678901234") // Use a securely generated key

	// Concatenate the parameters into a single string to encrypt
	parameterString := fmt.Sprintf("role=%s&userID=%d&adm=%s&username=%s&phone=%s&fee=%f",
		role, userID, adm, username, phone, fee)

	// Encrypt the parameter string
	encryptedParams, err := encrypt(parameterString, key)
	if err != nil {
		log.Printf("Error encrypting parameters: %v", err)
		return ""
	}

	return encryptedParams
}
