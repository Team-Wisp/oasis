package handler

// Handler for incoming HTTP POST req (email & OTP)

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/service"
)

type VerifyOTPRequest struct {
	EmailHash string `json:"email"`
	OTP       string `json:"otp"`
}

type VerifyOTPResponse struct {
	Verified bool `json:"verified"`
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyOTPRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	isValid := service.VerifyOTP(req.EmailHash, req.OTP)

	resp := VerifyOTPResponse{
		Verified: isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
