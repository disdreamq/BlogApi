package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/disdreamq/BlogApi/internal/port"
	"github.com/disdreamq/BlogApi/internal/service"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	userService port.UserService
}

type userRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(userService port.UserService) *UserController {
	return &UserController{userService: userService}

}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userReq userRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}
	user, err := c.userService.CreateUser(r.Context(), userReq.Username, userReq.Email, userReq.Password)
	if err != nil {
		switch err {
		case service.ErrUnexpected:
			http.Error(w, `{"error": "failed to create user"}`, http.StatusInternalServerError)
			return
		case service.ErrUserAlreadyExists:
			http.Error(w, `{"error": "user with this email already exists"}`, http.StatusConflict)
			return
		default:
			http.Error(w, `{"error": "failed to create user"}`, http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid userID"}`, http.StatusBadRequest)
		return
	}
	user, err := c.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "failed to get user"}`, http.StatusBadRequest)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	email := chi.URLParam(r, "user_email")
	user, err := c.userService.GetUserByEmail(r.Context(), email)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "failed to get user"}`, http.StatusBadRequest)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userReq userRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}
	err := c.userService.UpdateUser(r.Context(), userReq.Username, userReq.Email, userReq.Password)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		default:
			http.Error(w, `{"error": "failed to get user"}`, http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid userID"}`, http.StatusBadRequest)
		return
	}
	err = c.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "failed to get user"}`, http.StatusBadRequest)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
