package service

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const otpTTL = 5 * time.Minute

func hashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := sha256.Sum256([]byte(email))
	return fmt.Sprintf("%x", hash[:])
}

// GenerateAndStoreOTP creates a secure OTP and stores it in Redis with TTL
func GenerateAndStoreOTP(email string) (string, error) {
	code, err := generateOTP()
	if err != nil {
		return "", err
	}

	hashedEmail := hashEmail(email)
	key := fmt.Sprintf("otp:%s", hashedEmail)
	err = Redis.Set(Ctx, key, code, otpTTL).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifyOTP checks if the OTP is valid for the given email and deletes it if successful
func VerifyOTP(hashedEmail, otp string) bool {
	key := fmt.Sprintf("otp:%s", hashedEmail)

	storedCode, err := Redis.Get(Ctx, key).Result()
	if err == redis.Nil || err != nil {
		return false
	}

	if !subtleCompare(storedCode, otp) {
		return false
	}

	// Delete OTP after successful use
	Redis.Del(Ctx, key)
	return true
}

// generateOTP creates a secure random 6-digit OTP
func generateOTP() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

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
