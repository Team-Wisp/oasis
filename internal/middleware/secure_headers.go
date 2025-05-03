package middleware

import (
	"net/http"
	"os"
)

func SecureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("DESERT_URI") // for frame-ancestors

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' https://cdn.tailwindcss.com/3.4.16; frame-ancestors "+allowedOrigin+";")

		w.Header().Set("X-Frame-Options", "ALLOW-FROM "+allowedOrigin) // legacy
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=()")

		next.ServeHTTP(w, r)
	})
}
