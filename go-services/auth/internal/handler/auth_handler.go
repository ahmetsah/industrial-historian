package handler

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ahmet/historian/go-services/auth/internal/repository"
	"github.com/ahmet/historian/go-services/auth/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

const (
	SubjectLoginSuccess    = "sys.auth.login"
	SubjectSignatureIssued = "sys.auth.signature_issued"
)

type AuthHandler struct {
	Repo         repository.UserRepository
	TokenService *service.TokenService
	NatsConn     *nats.Conn
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.TokenService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Publish event to NATS
	event := map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Username,
		"timestamp": time.Now(),
		"event":     "login_success",
	}
	eventBytes, _ := json.Marshal(event)
	if h.NatsConn != nil {
		h.NatsConn.Publish(SubjectLoginSuccess, eventBytes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

type ReAuthRequest struct {
	Password string `json:"password"`
}

type ReAuthResponse struct {
	SigningToken string `json:"signing_token"`
}

func (h *AuthHandler) ReAuth(w http.ResponseWriter, r *http.Request) {
	// 1. Get Bearer Token
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	tokenString := authHeader[7:]

	// 2. Validate Token
	token, err := h.TokenService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims: username missing", http.StatusUnauthorized)
		return
	}

	// 3. Parse Request Body
	var req ReAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Verify Password
	user, err := h.Repo.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 5. Generate Signing Token
	signingToken, err := h.TokenService.GenerateSigningToken(user)
	if err != nil {
		http.Error(w, "Failed to generate signing token", http.StatusInternalServerError)
		return
	}

	// 6. Publish Event
	event := map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Username,
		"timestamp": time.Now(),
		"event":     "signature_issued",
	}
	eventBytes, _ := json.Marshal(event)
	if h.NatsConn != nil {
		h.NatsConn.Publish(SubjectSignatureIssued, eventBytes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ReAuthResponse{SigningToken: signingToken})
}

type CreateServiceAccountRequest struct {
	Name string `json:"name"`
}

type CreateServiceAccountResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) CreateServiceAccount(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Request
	var req CreateServiceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// 2. Create User with SERVICE role
	// Generate a random password (not used for service accounts, but required by DB)
	randomPass := make([]byte, 32)
	rand.Read(randomPass)
	passHash, _ := bcrypt.GenerateFromPassword(randomPass, bcrypt.DefaultCost)

	user := &repository.User{
		Username:     "svc_" + req.Name,
		PasswordHash: string(passHash),
		Role:         "SERVICE",
	}

	if err := h.Repo.CreateUser(user); err != nil {
		http.Error(w, "Failed to create service account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Generate Long-Lived Token
	// We need to extend TokenService to support custom expiration or just use GenerateToken
	// For MVP, we'll just use GenerateToken which gives 1 hour.
	// TODO: Update TokenService to support custom expiration for Service Accounts
	token, err := h.TokenService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateServiceAccountResponse{Token: token})
}
