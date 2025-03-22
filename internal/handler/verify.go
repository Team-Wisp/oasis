package handler

import (
	"fmt"
	"log"
	"net/http"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.Path)
	log.Printf("Headers: %v", r.Header)

	fmt.Fprintln(w, "✅ /verify route hit successfully")
}
