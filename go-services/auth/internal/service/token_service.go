package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenService(keyPath string) (*TokenService, error) {
	// Try to load key
	privKey, err := loadPrivateKey(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Generate new key
			privKey, err = generatePrivateKey(keyPath)
			if err != nil {
				return nil, fmt.Errorf("failed to generate private key: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}
	}

	return &TokenService{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing RSA private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func generatePrivateKey(path string) (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := pem.Encode(file, pemBlock); err != nil {
		return nil, err
	}

	return privKey, nil
}

func (s *TokenService) GenerateToken(user *repository.User) (string, error) {
	expiration := time.Hour * 1
	if user.Role == "SERVICE" {
		expiration = time.Hour * 24 * 365 // 1 year
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"type":     "access",
		"exp":      time.Now().Add(expiration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *TokenService) GenerateSigningToken(user *repository.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"type":     "signing",
		"scope":    "signature",
		"exp":      time.Now().Add(time.Minute * 1).Unix(), // 1 minute expiration
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *TokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
}
