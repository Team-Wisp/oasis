package service

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

const otpTTL = 5 * time.Minute
const maxFailedAttempts = 5
const blockDuration = 15 * time.Minute // Block for 15 minutes after exceeding the threshold

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
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	return code, nil
}

// VerifyOTP with RateLimit checks if the OTP is valid for the given email and deletes it if successful
func VerifyOTPWithRateLimit(hashedEmail, otp string) bool {
	key := fmt.Sprintf("otp:%s", hashedEmail)
	failedAttemptsKey := fmt.Sprintf("otp:failed:%s", hashedEmail)

	// Check if the user is blocked
	failedAttempts, _ := Redis.Get(Ctx, failedAttemptsKey).Int()
	if failedAttempts >= maxFailedAttempts {
		logger.WithFields(logrus.Fields{
			"hashedEmail":    hashedEmail,
			"failedAttempts": failedAttempts,
			"event":          "rate_limit_block",
		}).Warn("Too many failed OTP attempts. User is blocked.")
		fmt.Println("Too many failed attempts. Try again later.")
		return false
	}

	// Verify OTP
	storedCode, err := Redis.Get(Ctx, key).Result()
	if err == redis.Nil {
		logger.WithFields(logrus.Fields{
			"hashedEmail": hashedEmail,
			"event":       "otp_verification_failed",
			"reason":      "OTP not found or expired",
		}).Warn("Failed OTP verification")
		Redis.Incr(Ctx, failedAttemptsKey)                  // Increment failed attempts
		Redis.Expire(Ctx, failedAttemptsKey, blockDuration) // Set TTL for the block
		return false
	} else if err != nil {
		logger.WithFields(logrus.Fields{
			"hashedEmail": hashedEmail,
			"event":       "redis_error",
			"error":       err.Error(),
		}).Error("Redis error during OTP verification")
		return false
	}
	if !subtleCompare(storedCode, otp) {
		logger.WithFields(logrus.Fields{
			"hashedEmail": hashedEmail,
			"event":       "otp_verification_failed",
			"reason":      "Invalid OTP",
		}).Warn("Failed OTP verification")
		Redis.Incr(Ctx, failedAttemptsKey)                  // Increment failed attempts
		Redis.Expire(Ctx, failedAttemptsKey, blockDuration) // Set TTL for the block
		return false
	}

	logger.WithFields(logrus.Fields{
		"hashedEmail": hashedEmail,
		"event":       "otp_verification_success",
	}).Info("OTP verified successfully")

	// Delete OTP and reset failed attempts on success
	Redis.Del(Ctx, key)
	Redis.Del(Ctx, failedAttemptsKey)
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
