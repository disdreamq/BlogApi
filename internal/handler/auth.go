package handler

import (
	"encoding/json"
	"net/http"

	"github.com/disdreamq/BlogApi/internal/port"
	"github.com/disdreamq/BlogApi/internal/service"
)

type AuthController struct {
	authService port.AuthService
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthController(authService port.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var authReq AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}
	token, err := a.authService.Login(r.Context(), authReq.Email, authReq.Password)
	if err != nil {
		switch err {
		case service.ErrWrongPassword:
			http.Error(w, `{"error": "wrong password"}`, http.StatusUnauthorized)
			return
		case service.ErrCanNotLogin:
			http.Error(w, `{"error": "can not login"}`, http.StatusUnauthorized)
			return
		case service.ErrUserNotFound:
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		default:
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}
