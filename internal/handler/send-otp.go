package handler

// Handler for incoming HTTP POST req (email)

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Team-Wisp/oasis/internal/service"
)

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Message string `json:"message"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func SendOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Optional: Rate limiting here
	// if blocked := service.CheckAndThrottle(req.Email); blocked { ... }

	otp, err := service.GenerateAndStoreOTP(req.Email)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	subject := "Your Wisp Verification Code"
	bodyText := fmt.Sprintf("Hi,\n\nYour OTP is: %s\nIt expires in 5 minutes.\n\n– TeamWisp", otp)

	if err := service.SendEmail(req.Email, subject, bodyText); err != nil {
		http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SendOTPResponse{Message: "OTP sent to your email."})
}
