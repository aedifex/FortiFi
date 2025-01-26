package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),10)	
	if err != nil {
		return "", fmt.Errorf("could not create hash from password: %v", err.Error())
	}
	return string(hash), nil
}

func ValidatePassword(stored string, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	return err == nil
	
}