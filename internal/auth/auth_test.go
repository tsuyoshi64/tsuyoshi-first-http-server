package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTFlow(t *testing.T) {
	secret := "super-secret-key-12345"
	userID := uuid.New()

	// 1. Success Path
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Expected no error making JWT, got: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Expected valid verification token state, got: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected matching user ID %v, got %v", userID, parsedID)
	}

	// 2. Failure Path: Expired Token Rejection
	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Expected no error making expired JWT, got: %v", err)
	}

	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("Expected an validation error for an expired token payload, got nil")
	}

	// 3. Failure Path: Wrong Secret Rejection
	_, err = ValidateJWT(token, "completely-different-wrong-secret")
	if err == nil {
		t.Error("Expected validation error when using a compromised signature key, got nil")
	}
}
