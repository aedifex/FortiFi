package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

// Returns a JWT/Refresh Token pair.

// Returns non-nil error if unsuccessful.
func GenTokenPair(key string, id string) (string, *db.RefreshToken, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "FortiFi",
		Subject: id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute*10)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", nil, err
	}

	refreshString, err := genRefreshString()
	if err != nil {
		return "", nil, err
	}

	// Create refresh token with expiration date a week from now
	refresh := &db.RefreshToken{
		Id: id,
		Token: refreshString,
		Expires: SerializeTime(time.Now().Add(time.Hour*24*7)),
	}
	return signedToken, refresh, nil
}

// Validates the Jwt Signing algorithm and exp time then returns the associated subject id 
func GetJwtSubject(key string, signedToken string) (string, error) {
	
	// parse token and check signing method
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(key), nil
    }, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
    if err != nil {
        return "", fmt.Errorf("failed to parse token: %w", err)
    }

	// Check expiration date
	if expAt, err := parsedToken.Claims.GetExpirationTime(); (err != nil) || expAt.Time.Before(time.Now()){
		return "", errors.New("token has expired")
	}

	// Check signature
	if !parsedToken.Valid {
		return "", errors.New("invalid signature")
	}

	sub, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("failed to get claims subject: %s", err.Error())
	}
	return sub, nil
}

func genRefreshString() (string, error) {
	// Generate random 20 byte string
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(bytes)

	return tokenString, nil
}

func ExtractBearer(tokenHeader string) (string, error) {
	signedToken := ""
	if strings.HasPrefix(tokenHeader, "Bearer ") {
		signedToken = strings.TrimPrefix(tokenHeader, "Bearer ")
	} else {
		return "", errors.New("invalid Authorization header")
	}
	return signedToken, nil
}