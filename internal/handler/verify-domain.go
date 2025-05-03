package handler

// Handler for incoming HTTP POST req (email domains)
import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Team-Wisp/oasis/internal/service"
)

type VerifyRequest struct {
	Email string `json:"email"`
}

type VerifyResponse struct {
	Domain  string `json:"domain"`
	IsValid bool   `json:"isValid"`
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, "Valid email is required", http.StatusBadRequest)
		return
	}
	domain := strings.Split(email, "@")[1]

	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	isValid := service.CheckMX(domain)
	domainType := service.GetDomainType(domain)

	if domainType == "public" {
		http.Error(w, "Public email domains are not allowed, only college/corporate email!!", http.StatusForbidden)
		return
	}

	// Enrich org info in background
	if isValid {
		service.GetOrInitDomain(domain, domainType)
	}

	resp := VerifyResponse{
		Domain:  domain,
		IsValid: isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
