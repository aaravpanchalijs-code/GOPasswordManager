package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func Encrypt(plainText string) (string, error) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		return "", errors.New("ENCRYPTION_KEY not found")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(
		nonce,
		nonce,
		[]byte(plainText),
		nil,
	)
	return base64.StdEncoding.EncodeToString(cipherText), nil

}

func Decrypt(cipherText string) (string, error) {

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		return "", errors.New("ENCRYPTION_KEY not found")
	}

	// Convert Base64 string back to bytes
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Create AES-GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get nonce size
	nonceSize := aesGCM.NonceSize()

	// Validate ciphertext length
	if len(cipherBytes) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	// Separate nonce and encrypted data
	nonce := cipherBytes[:nonceSize]
	cipherBytes = cipherBytes[nonceSize:]

	// Decrypt
	plainText, err := aesGCM.Open(
		nil,
		nonce,
		cipherBytes,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
