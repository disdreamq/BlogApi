package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/service"
)

type mockUserService struct {
	mockCreateFunc     func(ctx context.Context, username, email, password string) (*domain.User, error)
	mockGetByIDFunc    func(ctx context.Context, userID int64) (*domain.User, error)
	mockGetByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	mockUpdateFunc     func(ctx context.Context, currUserID, userID int64, username, email, password string) error
	mockDeleteFunc     func(ctx context.Context, currUserID, userID int64) error
}

func (m *mockUserService) Create(ctx context.Context, username, email, password string) (*domain.User, error) {
	return m.mockCreateFunc(ctx, username, email, password)
}
func (m *mockUserService) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	return m.mockGetByIDFunc(ctx, userID)
}
func (m *mockUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.mockGetByEmailFunc(ctx, email)
}
func (m *mockUserService) Update(ctx context.Context, currUserID, userID int64, username, email, password string) error {
	return m.mockUpdateFunc(ctx, currUserID, userID, username, email, password)
}
func (m *mockUserService) Delete(ctx context.Context, currUserID, userID int64) error {
	return m.mockDeleteFunc(ctx, currUserID, userID)
}

func TestUserController_Create(t *testing.T) {
	createdUser := userResponse{ID: 67, Username: "johndoe", Email: "user@example.com"}
	user, _ := json.Marshal(createdUser)
	resp := string(user)
	tests := []struct {
		name           string
		service        mockUserService
		input          createUserRequest
		expectedStatus int
		expectedBody   string
	}{
		{"happy path", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return &domain.User{ID: 67, Username: username, Email: email, PasswordHash: password}, nil
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "password123"}, http.StatusCreated, resp},
		{"empty username", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidUserName
		}}, createUserRequest{Username: "", Email: "user@example.com", Password: "password123"}, http.StatusBadRequest, ""},
		{"empty email", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidEmail
		}}, createUserRequest{Username: "johndoe", Email: "", Password: "password123"}, http.StatusBadRequest, ""},
		{"empty password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: ""}, http.StatusBadRequest, ""},
		{"invalid email", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidEmail
		}}, createUserRequest{Username: "johndoe", Email: "invalidemail", Password: "password123"}, http.StatusBadRequest, ""},
		{"too long username", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidUserName
		}}, createUserRequest{Username: strings.Repeat("67", 31), Email: "user@example.com", Password: "password123"}, http.StatusBadRequest, ""},
		{"too long password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: strings.Repeat("67", 61)}, http.StatusBadRequest, ""},
		{"too short password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "short"}, http.StatusBadRequest, ""},
		{"user already exists", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrUserAlreadyExists
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "short"}, http.StatusConflict, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewUserController(&tt.service)
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctrl.Create(w, req)
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if tt.expectedBody != "" {
				var got, want interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.expectedBody), &want); err != nil {
					t.Errorf("failed to unmarshal expected: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("expected body %v, got %v", want, got)
				}
			}
		})
	}
}

func TestUserController_GetByID(t *testing.T) {
	prepUser := userResponse{ID: 1, Username: "johndoe", Email: "user@example.com"}
	user, _ := json.Marshal(prepUser)
	resp := string(user)
	tests := []struct {
		name           string
		service        mockUserService
		input          createUserRequest
		expectedStatus int
		expectedBody   string
	}{
		{"happy path", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: username, Email: email, PasswordHash: password}, nil
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "password123"}, http.StatusCreated, resp},
		{"empty username", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidUserName
		}}, createUserRequest{Username: "", Email: "user@example.com", Password: "password123"}, http.StatusBadRequest, ""},
		{"empty email", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidEmail
		}}, createUserRequest{Username: "johndoe", Email: "", Password: "password123"}, http.StatusBadRequest, ""},
		{"empty password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: ""}, http.StatusBadRequest, ""},
		{"invalid email", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidEmail
		}}, createUserRequest{Username: "johndoe", Email: "invalidemail", Password: "password123"}, http.StatusBadRequest, ""},
		{"too long username", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, domain.ErrInvalidUserName
		}}, createUserRequest{Username: strings.Repeat("67", 31), Email: "user@example.com", Password: "password123"}, http.StatusBadRequest, ""},
		{"too long password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: strings.Repeat("67", 61)}, http.StatusBadRequest, ""},
		{"too short password", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrInvalidPasswordLength
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "short"}, http.StatusBadRequest, ""},
		{"user already exists", mockUserService{mockCreateFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
			return nil, service.ErrUserAlreadyExists
		}}, createUserRequest{Username: "johndoe", Email: "user@example.com", Password: "short"}, http.StatusConflict, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewUserController(&tt.service)
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctrl.Create(w, req)
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if tt.expectedBody != "" {
				var got, want interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.expectedBody), &want); err != nil {
					t.Errorf("failed to unmarshal expected: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("expected body %v, got %v", want, got)
				}
			}
		})
	}
}
