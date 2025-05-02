package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/service"
)

type LoginRequest struct {
	Email    string `json:"email"`    // hashed
	Password string `json:"password"` // hashed
}

type LoginResponse struct {
	Token string `json:"token"`
}

func VerifyLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if !service.IsValidEmailHash(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	user, err := service.GetUserByEmailHash(req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !service.CheckPassword(user.Password, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := service.GenerateJWT(req.Email, user.OrgSlug, user.OrgType)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
