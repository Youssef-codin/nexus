package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Youssef-codin/NexusPay/internal/security"
)

const testRefreshTokenDuration = 7 * 24 * time.Hour

type mockAuthService struct {
	loginFunc    func(ctx context.Context, req loginRequest) (loginResponse, error)
	registerFunc func(ctx context.Context, req registerRequest) (registerResponse, error)
	refreshFunc  func(ctx context.Context, req refreshRequest) (refreshResponse, error)
	logoutFunc   func(ctx context.Context) error
}

func (m *mockAuthService) login(ctx context.Context, req loginRequest) (loginResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return loginResponse{}, nil
}

func (m *mockAuthService) register(ctx context.Context, req registerRequest) (registerResponse, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return registerResponse{}, nil
}

func (m *mockAuthService) refreshToken(ctx context.Context, req refreshRequest) (refreshResponse, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, req)
	}
	return refreshResponse{}, nil
}

func (m *mockAuthService) logout(ctx context.Context) error {
	if m.logoutFunc != nil {
		return m.logoutFunc(ctx)
	}
	return nil
}

func TestLoginController_Success(t *testing.T) {
	mock := &mockAuthService{
		loginFunc: func(ctx context.Context, req loginRequest) (loginResponse, error) {
			return loginResponse{
				Email:        req.Email,
				FullName:     "Test User",
				JwtToken:     "jwt-token",
				RefreshToken: "refresh-token",
			}, nil
		},
	}

	ctrl := NewController(mock)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.LoginController(w, req)
	if err != nil {
		t.Fatalf("controller returned error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp loginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.Email)
	}
	if resp.JwtToken == "" {
		t.Error("expected JWT token in response")
	}
}

func TestLoginController_BadRequest(t *testing.T) {
	mock := &mockAuthService{}
	ctrl := NewController(mock)

	body := `{"email":"invalid-email"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.LoginController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", we.StatusCode())
	}
}

func TestLoginController_Unauthorized(t *testing.T) {
	mock := &mockAuthService{
		loginFunc: func(ctx context.Context, req loginRequest) (loginResponse, error) {
			return loginResponse{}, ErrInvalidCredentials
		},
	}

	ctrl := NewController(mock)

	body := `{"email":"test@example.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.LoginController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", we.StatusCode())
	}
}

func TestLoginController_UserNotFound(t *testing.T) {
	mock := &mockAuthService{
		loginFunc: func(ctx context.Context, req loginRequest) (loginResponse, error) {
			return loginResponse{}, ErrUserNotFound
		},
	}

	ctrl := NewController(mock)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.LoginController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", we.StatusCode())
	}
}

func TestRegisterController_Success(t *testing.T) {
	mock := &mockAuthService{
		registerFunc: func(ctx context.Context, req registerRequest) (registerResponse, error) {
			return registerResponse{
				Email:        req.Email,
				FullName:     req.FullName,
				JwtToken:     "jwt-token",
				RefreshToken: "refresh-token",
			}, nil
		},
	}

	ctrl := NewController(mock)

	body := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RegisterController(w, req)
	if err != nil {
		t.Fatalf("controller returned error: %v", err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp registerResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.Email)
	}
}

func TestRegisterController_BadRequest(t *testing.T) {
	mock := &mockAuthService{}
	ctrl := NewController(mock)

	body := `{"email":"invalid-email"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RegisterController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", we.StatusCode())
	}
}

func TestRegisterController_UserAlreadyExists(t *testing.T) {
	mock := &mockAuthService{
		registerFunc: func(ctx context.Context, req registerRequest) (registerResponse, error) {
			return registerResponse{}, ErrUserAlreadyExists
		},
	}

	ctrl := NewController(mock)

	body := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RegisterController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusConflict {
		t.Errorf("expected status 409, got %d", we.StatusCode())
	}
}

func TestRefreshController_Success(t *testing.T) {
	mock := &mockAuthService{
		refreshFunc: func(ctx context.Context, req refreshRequest) (refreshResponse, error) {
			return refreshResponse{
				JwtToken:     "new-jwt-token",
				RefreshToken: "new-refresh-token",
			}, nil
		},
	}

	ctrl := NewController(mock)

	body := `{"refresh_token":"old-token"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RefreshController(w, req)
	if err != nil {
		t.Fatalf("controller returned error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp refreshResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.JwtToken == "" {
		t.Error("expected JWT token in response")
	}
}

func TestRefreshController_BadRequest(t *testing.T) {
	mock := &mockAuthService{}
	ctrl := NewController(mock)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RefreshController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", we.StatusCode())
	}
}

func TestRefreshController_UserNotFound(t *testing.T) {
	mock := &mockAuthService{
		refreshFunc: func(ctx context.Context, req refreshRequest) (refreshResponse, error) {
			return refreshResponse{}, ErrUserNotFound
		},
	}

	ctrl := NewController(mock)

	body := `{"refresh_token":"invalid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RefreshController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", we.StatusCode())
	}
}

func TestRefreshController_TokenExpired(t *testing.T) {
	mock := &mockAuthService{
		refreshFunc: func(ctx context.Context, req refreshRequest) (refreshResponse, error) {
			return refreshResponse{}, ErrTokenExpired
		},
	}

	ctrl := NewController(mock)

	body := `{"refresh_token":"expired-token"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := ctrl.RefreshController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", we.StatusCode())
	}
}

func TestTestAuth(t *testing.T) {
	auth := security.NewAuthenticator("test-secret", testRefreshTokenDuration)
	ctrl := &controller{svc: nil}

	router := http.NewServeMux()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		ctrl.TestAuth(w, r)
	})

	token := auth.MakeJWTToken(security.Claims{ID: "user123"})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestTestAuth_NoToken(t *testing.T) {
	ctrl := &controller{svc: nil}

	router := http.NewServeMux()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		ctrl.TestAuth(w, r)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestTestAuth_InvalidToken(t *testing.T) {
	ctrl := &controller{svc: nil}

	router := http.NewServeMux()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		ctrl.TestAuth(w, r)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestLogoutController_Success(t *testing.T) {
	mock := &mockAuthService{
		logoutFunc: func(ctx context.Context) error {
			return nil
		},
	}

	ctrl := NewController(mock)

	auth := security.NewAuthenticator("test-secret", testRefreshTokenDuration)
	token := auth.MakeJWTToken(security.Claims{ID: "user123"})

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	err := ctrl.LogoutController(w, req)
	if err != nil {
		t.Fatalf("controller returned error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestLogoutController_UserNotFound(t *testing.T) {
	mock := &mockAuthService{
		logoutFunc: func(ctx context.Context) error {
			return ErrUserNotFound
		},
	}

	ctrl := NewController(mock)

	auth := security.NewAuthenticator("test-secret", testRefreshTokenDuration)
	token := auth.MakeJWTToken(security.Claims{ID: "user123"})

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	err := ctrl.LogoutController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	we, ok := err.(interface{ StatusCode() int })
	if !ok {
		t.Fatalf("expected WrappedError, got %T", err)
	}
	if we.StatusCode() != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", we.StatusCode())
	}
}

type loggedError struct {
	msg string
}

func (e loggedError) Error() string {
	return e.msg
}

func TestLogoutController_LogsError(t *testing.T) {
	mock := &mockAuthService{
		logoutFunc: func(ctx context.Context) error {
			return errors.New("some error")
		},
	}

	ctrl := NewController(mock)

	auth := security.NewAuthenticator("test-secret", testRefreshTokenDuration)
	token := auth.MakeJWTToken(security.Claims{ID: "user123"})

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	err := ctrl.LogoutController(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
