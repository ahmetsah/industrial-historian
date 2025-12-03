package service

import (
	"os"
	"testing"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	keyPath := "test_private.pem"
	defer os.Remove(keyPath)

	service, err := NewTokenService(keyPath)
	if err != nil {
		t.Fatalf("Failed to create token service: %v", err)
	}

	user := &repository.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
	}

	tokenString, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Token string is empty")
	}

	// Verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.publicKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatal("Token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims")
	}

	if claims["username"] != "testuser" {
		t.Errorf("Expected username testuser, got %v", claims["username"])
	}
	if claims["role"] != "admin" {
		t.Errorf("Expected role admin, got %v", claims["role"])
	}
}

func TestGenerateSigningToken(t *testing.T) {
	keyPath := "test_private_signing.pem"
	defer os.Remove(keyPath)

	service, err := NewTokenService(keyPath)
	if err != nil {
		t.Fatalf("Failed to create token service: %v", err)
	}

	user := &repository.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
	}

	tokenString, err := service.GenerateSigningToken(user)
	if err != nil {
		t.Fatalf("Failed to generate signing token: %v", err)
	}

	// Verify token
	token, err := service.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims")
	}

	if claims["type"] != "signing" {
		t.Errorf("Expected type signing, got %v", claims["type"])
	}
	if claims["scope"] != "signature" {
		t.Errorf("Expected scope signature, got %v", claims["scope"])
	}
}
