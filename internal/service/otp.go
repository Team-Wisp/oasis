package service

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

type otpEntry struct {
	Code      string
	ExpiresAt time.Time
}

var (
	otpStore = make(map[string]otpEntry)
	otpLock  sync.Mutex
	otpTTL   = 5 * time.Minute
)

// GenerateAndStoreOTP creates a new OTP, stores it in memory, and returns the OTP
func GenerateAndStoreOTP(email string) (string, error) {
	code, err := generateOTP()
	if err != nil {
		return "", err
	}

	otpLock.Lock()
	defer otpLock.Unlock()

	otpStore[email] = otpEntry{
		Code:      code,
		ExpiresAt: time.Now().Add(otpTTL),
	}

	return code, nil
}

// VerifyOTP checks if the OTP is valid and not expired
func VerifyOTP(email string, otp string) bool {
	otpLock.Lock()
	defer otpLock.Unlock()

	entry, exists := otpStore[email]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return false
	}

	// Optionally: delete the OTP after successful use
	delete(otpStore, email)

	// Constant-time comparison for security
	return subtleCompare(entry.Code, otp)
}

// generateOTP creates a secure random 6-digit OTP
func generateOTP() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to 6-digit number
	num := (int(b[0]) << 16) | (int(b[1]) << 8) | int(b[2])
	code := fmt.Sprintf("%06d", num%1000000)
	return code, nil
}

// subtleCompare prevents timing attacks
func subtleCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i] ^ b[i])
	}
	return result == 0
}

// CleanupExpiredOTPs removes expired OTPs (optional for a background goroutine)
func CleanupExpiredOTPs() {
	otpLock.Lock()
	defer otpLock.Unlock()

	now := time.Now()
	for email, entry := range otpStore {
		if now.After(entry.ExpiresAt) {
			delete(otpStore, email)
		}
	}
}
