package main

import (
	"log"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/handler"
	"github.com/Team-Wisp/oasis/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}
	//connect to db
	service.InitRedis()
	service.InitMongo()

	http.HandleFunc("/verify", handler.VerifyHandler)        // server/verfiy
	http.HandleFunc("/send-otp", handler.SendOTPHandler)     // server/send-otp
	http.HandleFunc("/verify-otp", handler.VerifyOTPHandler) // server/verify-otp

	log.Println("🌿 Oasis is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
