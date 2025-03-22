package main

import (
	"log"
	"net/http"

	"github.com/Team-Wisp/oasis/internal/handler"
)

func main() {
	http.HandleFunc("/verify", handler.VerifyHandler)

	log.Println("🌿 Oasis is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
