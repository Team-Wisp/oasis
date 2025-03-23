package service

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const otpTTL = 5 * time.Minute

// GenerateAndStoreOTP creates a secure OTP and stores it in Redis with TTL
func GenerateAndStoreOTP(email string) (string, error) {
	code, err := generateOTP()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("otp:%s", email)
	err = Redis.Set(Ctx, key, code, otpTTL).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifyOTP checks if the OTP is valid for the given email and deletes it if successful
func VerifyOTP(email, otp string) bool {
	key := fmt.Sprintf("otp:%s", email)

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
