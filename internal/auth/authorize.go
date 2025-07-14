package auth

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {	
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if tokenSecret == "" {
		return "", fmt.Errorf("tokenSecret must not be blank")
	}

	claims := jwt.RegisteredClaims{
		Issuer:	"chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error making token string: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func (*jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error: token is malformed, expired, or tampered -- %w", err)
	} else if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("Error: token is invalid")
	}

	if claims.Subject == "" {
		return uuid.UUID{}, fmt.Errorf("Error: no subject")
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing UUID: %w", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headerInfo := headers.Get("Authorization")
	if headerInfo == "" {
		return "", fmt.Errorf("Error getting header information")
	}
	headerParts := strings.Split(headerInfo, " ")
	if len(headerParts) <= 1 {
		return "", fmt.Errorf("Error getting header information")
	}

	return headerParts[1], nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	headerInfo := headers.Get("Authorization")
	if headerInfo == "" {
		return "", fmt.Errorf("Errorr getting headerr information")
	}
	headerParts := strings.Split(headerInfo, " ")
	if len(headerParts) <= 1 {
		return "", fmt.Errorf("Error getting header information")
	}

	return headerParts[1], nil
}
