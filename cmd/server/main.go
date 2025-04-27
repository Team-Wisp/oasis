package main

import (
	"log"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/handler"
	"github.com/Team-Wisp/oasis/internal/service"
	"github.com/joho/godotenv"
)

// CORS middleware - browser security feature & does not apply to server-to-server HTTP requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow Desert frontend during local dev
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
		return
	}
	//connect to db
	service.InitRedis()
	service.InitMongo()

	mux := http.NewServeMux()

	mux.Handle("/verify-domain", corsMiddleware(http.HandlerFunc(handler.VerifyHandler)))
	mux.Handle("/send-otp", corsMiddleware(http.HandlerFunc(handler.SendOTPHandler)))
	mux.Handle("/verify-otp", corsMiddleware(http.HandlerFunc(handler.VerifyOTPHandler)))

	log.Println("🌿 Oasis is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
