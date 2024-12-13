package handlers

import (
	"crypto/aes"
	"crypto/cipher"
"fmt"
	"encoding/base64"
	
)


// Decrypt function
func decrypt(ciphertext string, key []byte) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Get the IV from the ciphertext
	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	// Decrypt the data
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	// Return the decrypted string
	return string(data), nil
}
