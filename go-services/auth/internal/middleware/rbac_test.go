package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/repository"
	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/service"
)

func TestRBACMiddleware(t *testing.T) {
	// Setup
	tokenService, _ := service.NewTokenService("test_private_rbac.pem")
	defer os.Remove("test_private_rbac.pem")

	rbac := NewRBACMiddleware(tokenService)

	// Create tokens
	adminUser := &repository.User{ID: 1, Username: "admin", Role: "ADMIN"}
	adminToken, _ := tokenService.GenerateToken(adminUser)

	operatorUser := &repository.User{ID: 2, Username: "operator", Role: "OPERATOR"}
	operatorToken, _ := tokenService.GenerateToken(operatorUser)

	// Test Handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 1. Test Admin Access to Admin Route
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()

	rbac.Authenticate(rbac.RequireRole("ADMIN")(testHandler)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for Admin, got %d", w.Code)
	}

	// 2. Test Operator Access to Admin Route
	reqOp := httptest.NewRequest("GET", "/", nil)
	reqOp.Header.Set("Authorization", "Bearer "+operatorToken)
	wOp := httptest.NewRecorder()

	rbac.Authenticate(rbac.RequireRole("ADMIN")(testHandler)).ServeHTTP(wOp, reqOp)

	if wOp.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden for Operator, got %d", wOp.Code)
	}

	// 3. Test No Token
	reqNoToken := httptest.NewRequest("GET", "/", nil)
	wNoToken := httptest.NewRecorder()

	rbac.Authenticate(rbac.RequireRole("ADMIN")(testHandler)).ServeHTTP(wNoToken, reqNoToken)

	if wNoToken.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for no token, got %d", wNoToken.Code)
	}
}
