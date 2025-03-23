package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/service"
)

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Message string `json:"message"`
}

func SendOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendOTPRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Generate and store OTP
	otp, err := service.GenerateAndStoreOTP(req.Email)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	// Email content
	subject := "Your TeamWisp Verification Code"
	bodyText := fmt.Sprintf(`Hi,

Your One-Time Password (OTP) is: %s

It will expire in 5 minutes.

– TeamWisp`, otp)

	// Send email
	if err := service.SendEmail(req.Email, subject, bodyText); err != nil {
		http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SendOTPResponse{
		Message: "OTP sent to your email.",
	})
}
