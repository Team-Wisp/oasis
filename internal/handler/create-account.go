package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Team-Wisp/oasis/internal/service"
)

type CreateAccountRequest struct {
	Email    string `json:"email"`    // SHA256
	Password string `json:"password"` // SHA256
	Domain   string `json:"domain"`   // plain domain string
}

type CreateAccountResponse struct {
	Message string `json:"message"`
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	exists, err := service.DoesUserExist(req.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Validate domain
	if !isValidDomain(req.Domain) {
		http.Error(w, "Invalid domain format", http.StatusBadRequest)
		return
	}

	// Get organization info
	org, err := service.LookupOrg(req.Domain)
	if err != nil {
		http.Error(w, "Could not find organization info", http.StatusBadRequest)
		return
	}

	// Final password hash (bcrypt of client-side hashed password)
	bcryptHash, err := service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Store user
	user := service.User{
		EmailHash: req.Email,
		Password:  bcryptHash,
		OrgSlug:   org.OrgSlug,
		OrgType:   org.OrgType,
		CreatedAt: time.Now(),
	}

	if err := service.SaveUser(user); err != nil {
		log.Printf("SaveUser failed: %+v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CreateAccountResponse{Message: "User created successfully"})
}

func isValidDomain(domain string) bool {
	// Regular expression to validate domain names
	var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}
