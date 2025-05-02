package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Team-Wisp/oasis/internal/handler"
	"github.com/Team-Wisp/oasis/internal/middleware"
	"github.com/Team-Wisp/oasis/internal/service"
	"github.com/joho/godotenv"
)

// CORS middleware - browser security feature & does not apply to server-to-server HTTP requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("DESERT_URI")) // Allow communication b/w services
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

	// Serve static files from /static directory
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Auth & login pages with secure headers
	mux.Handle("/signup", middleware.SecureHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/signup.html")
	})))

	mux.Handle("/login", middleware.SecureHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})))

	// mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./static/signup.html")
	// })

	// mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./static/login.html")
	// })

	//API Endpoints
	mux.Handle("/verify-domain", corsMiddleware(http.HandlerFunc(handler.VerifyHandler)))
	mux.Handle("/send-otp", corsMiddleware(http.HandlerFunc(handler.SendOTPHandler)))
	mux.Handle("/verify-otp", corsMiddleware(http.HandlerFunc(handler.VerifyOTPHandler)))
	mux.Handle("/create-account", corsMiddleware(http.HandlerFunc(handler.CreateAccountHandler)))
	mux.Handle("/verify-login", corsMiddleware(http.HandlerFunc(handler.VerifyLoginHandler)))

	log.Println("🌿 Oasis is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
