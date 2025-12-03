package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/repository"
	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository for integration testing without DB
type MockUserRepository struct {
	users map[string]*repository.User
}

func (m *MockUserRepository) CreateUser(user *repository.User) error {
	if m.users == nil {
		m.users = make(map[string]*repository.User)
	}
	user.ID = int64(len(m.users) + 1)
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*repository.User, error) {
	if user, ok := m.users[username]; ok {
		return user, nil
	}
	return nil, nil
}

func TestLoginFlow(t *testing.T) {
	// Setup
	repo := &MockUserRepository{}
	
	// Create test user
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo.CreateUser(&repository.User{
		Username:     "testuser",
		PasswordHash: string(passHash),
		Role:         "user",
	})

	tokenService, _ := service.NewTokenService("test_private_integration.pem")
	defer os.Remove("test_private_integration.pem")

	// No NATS for this test
	h := &AuthHandler{
		Repo:         repo,
		TokenService: tokenService,
		NatsConn:     nil, 
	}

	// Test Login Success
	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	h.Login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	if loginResp.Token == "" {
		t.Error("Expected token, got empty string")
	}

	// Test Login Failure
	reqBodyFail, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})
	reqFail := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBodyFail))
	wFail := httptest.NewRecorder()

	h.Login(wFail, reqFail)

	respFail := wFail.Result()
	if respFail.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", respFail.StatusCode)
	}
}

func TestReAuthFlow(t *testing.T) {
	// Setup
	repo := &MockUserRepository{}
	
	// Create test user
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &repository.User{
		Username:     "testuser",
		PasswordHash: string(passHash),
		Role:         "user",
	}
	repo.CreateUser(user)

	tokenService, _ := service.NewTokenService("test_private_reauth.pem")
	defer os.Remove("test_private_reauth.pem")

	h := &AuthHandler{
		Repo:         repo,
		TokenService: tokenService,
		NatsConn:     nil, 
	}

	// 1. Login to get Access Token
	accessToken, _ := tokenService.GenerateToken(user)

	// 2. Test ReAuth Success
	reqBody, _ := json.Marshal(map[string]string{
		"password": "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/re-auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()

	h.ReAuth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var reAuthResp ReAuthResponse
	json.NewDecoder(resp.Body).Decode(&reAuthResp)
	if reAuthResp.SigningToken == "" {
		t.Error("Expected signing token, got empty string")
	}

	// 3. Test ReAuth Failure (Wrong Password)
	reqBodyFail, _ := json.Marshal(map[string]string{
		"password": "wrongpassword",
	})
	reqFail := httptest.NewRequest("POST", "/api/v1/re-auth", bytes.NewBuffer(reqBodyFail))
	reqFail.Header.Set("Authorization", "Bearer "+accessToken)
	wFail := httptest.NewRecorder()

	h.ReAuth(wFail, reqFail)

	respFail := wFail.Result()
	if respFail.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", respFail.StatusCode)
	}

	// 4. Test ReAuth Failure (Invalid Token)
	reqBodyInvalid, _ := json.Marshal(map[string]string{
		"password": "password123",
	})
	reqInvalid := httptest.NewRequest("POST", "/api/v1/re-auth", bytes.NewBuffer(reqBodyInvalid))
	reqInvalid.Header.Set("Authorization", "Bearer invalidtoken")
	wInvalid := httptest.NewRecorder()

	h.ReAuth(wInvalid, reqInvalid)

	respInvalid := wInvalid.Result()
	if respInvalid.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", respInvalid.StatusCode)
	}
}

func TestCreateServiceAccount(t *testing.T) {
	// Setup
	repo := &MockUserRepository{}
	tokenService, _ := service.NewTokenService("test_private_svc.pem")
	defer os.Remove("test_private_svc.pem")

	h := &AuthHandler{
		Repo:         repo,
		TokenService: tokenService,
		NatsConn:     nil,
	}

	// Test Success
	reqBody, _ := json.Marshal(map[string]string{
		"name": "ingestor",
	})
	req := httptest.NewRequest("POST", "/api/v1/service-accounts", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	h.CreateServiceAccount(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var svcResp CreateServiceAccountResponse
	json.NewDecoder(resp.Body).Decode(&svcResp)
	if svcResp.Token == "" {
		t.Error("Expected token, got empty string")
	}

	// Verify User Created
	user, _ := repo.GetUserByUsername("svc_ingestor")
	if user == nil {
		t.Error("Expected user svc_ingestor to be created")
	}
	if user.Role != "SERVICE" {
		t.Errorf("Expected role SERVICE, got %s", user.Role)
	}
}
