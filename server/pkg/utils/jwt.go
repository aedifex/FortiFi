package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenJwt(key string, userId string) (string, string, time.Time, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "FortiFi",
		Subject: userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour*24)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", "", time.Time{}, err
	}

	refresh, expTime, err := genRefresh(userId)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, refresh, expTime, nil
}

// Validates the Jwt Signing algorithm and exp time then returns the associated user id 
func GetJwtId(key string, tokenHeader string) (string,error) {

	signedToken := ""
	if strings.HasPrefix(tokenHeader, "Bearer ") {
		signedToken = strings.TrimPrefix(tokenHeader, "Bearer ")
	} else {
		return "", errors.New("invalid Authorization header")
	}

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
		return "", fmt.Errorf("invalid signature")
	}

	sub, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("failed to get claims subject: %s", err.Error())
	}
	return sub, nil
}

func genRefresh(userId string) (string, time.Time, error) {
	// Generate random 20 byte string
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", time.Time{}, err
	}
	key := hex.EncodeToString(bytes)

	// Generate the refresh token
	expAt := time.Now().Add(time.Hour*24*7)
	claims := jwt.RegisteredClaims{
		Issuer: "FortiFi",
		Subject: userId,
		ExpiresAt: jwt.NewNumericDate(expAt),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %s", err.Error())
	}

	return signedToken, expAt, nil
}