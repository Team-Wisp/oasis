package handler

import (
	"encoding/json"
	// "fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Team-Wisp/oasis/internal/service"
)

// request type
type VerifyRequest struct {
	Email string `json:"email"`
}

// response type
type VerifyResponse struct {
	Domain       string `json:"domain"`
	Organization string `json:"organization"`
	IsValid      bool   `json:"isValid"`
}

// /verify
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Error handler
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

	log.Printf("📩 Email Received: %s", req.Email)

	// Extract domain
	parts := strings.Split(req.Email, "@")
	if len(parts) != 2 {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	domain := parts[1] // ex:someone@solarwinds.com (extract solarwinds.com)

	// Verify domain using net.LookupMX
	isValid := service.CheckMX(domain)

	org := service.MapDomainToOrg(domain)

	resp := VerifyResponse{
		Domain:       domain,
		Organization: org,
		IsValid:      isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
