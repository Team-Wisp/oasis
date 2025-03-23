package main

import (
	"log"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}
	http.HandleFunc("/verify", handler.VerifyHandler)    // server/verfiy
	http.HandleFunc("/send-otp", handler.SendOTPHandler) // server/send-otp

	log.Println("🌿 Oasis is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
