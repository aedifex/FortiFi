package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(plaintext string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext),10)	
	if err != nil {
		return "", fmt.Errorf("could not create hash from password: %v", err.Error())
	}
	return string(hash), nil
}

func HashMatch(stored string, given string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(given))
	return err == nil
	
}