package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"net/http"
)

func TestSuccessfulValidation(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
		return
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Failed to validate: %v", err)
		return
	}

	if validatedID != userID {
		t.Errorf("Expected user ID: %v Validated ID: %v", userID, validatedID)
		return
	}
}

func TestExpiredTokens(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	expiresIn := time.Millisecond

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
		return
	}

	time.Sleep(time.Millisecond * 10)

	_, err = ValidateJWT(token, secret)
	if err != nil {
		return
	}

	t.Errorf("Expected failed validation, but it succeeded")
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	duration := time.Minute * 5

	tokenString, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
		return
	}

	_, err = ValidateJWT(tokenString, "decret")
	if err != nil {
		return
	}

	t.Error("Error: token should not have been valid")
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "bearer 12345")

	token, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Error getting bearer token")
	}

	if token != "12345" {
		t.Errorf("Token strings don't match")
	}
}
