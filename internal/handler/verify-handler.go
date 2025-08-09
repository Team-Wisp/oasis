package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
		http.Error(w, "Email Does not exist", http.StatusBadRequest)
		return
	}
	if !service.CheckPassword(user.Password, req.Password) {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	// Resolve org by stored slug
	org, err := service.LookupOrgBySlug(user.Slug)
	if err != nil {
		http.Error(w, "Organization not found", http.StatusUnauthorized)
		return
	}

	// Ensure membership exists
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mid, _, err := service.UpsertMembership(ctx, user.ID, org.ID)
	if err != nil {
		http.Error(w, "Failed to init membership", http.StatusInternalServerError)
		return
	}

	token, err := service.GenerateJWT(user.ID.Hex(), user.Slug, user.OrgType, mid.Hex())
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: token})

}
