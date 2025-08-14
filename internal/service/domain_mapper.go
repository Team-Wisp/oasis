package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var httpClient = &http.Client{
	Timeout: 6 * time.Second,
}

func GetOrInitDomain(domain, domainType string) {
	// routine, context + timeout
	go func(domain, domainType string) {
		ctx, cancel := context.WithTimeout(context.Background(), httpClient.Timeout)
		defer cancel()

		payload := map[string]string{
			"domain":     domain,
			"domainType": domainType, // "corporate" | "college"
		}
		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Failed to marshal enrichment payload:", err)
			return
		}

		url := os.Getenv("DESERT_ENRICH_URI") //api-end point for domain enrichment
		if url == "" {
			fmt.Println("DESERT_ENRICH_URL is not set")
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Failed to create HTTP request to Desert:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Shared-secret auth
		secret := os.Getenv("AUTH_WEBHOOK_SECRET")
		if secret == "" {
			fmt.Println("AUTH_WEBHOOK_SECRET not set")
			return
		}
		req.Header.Set("X-Auth-Secret", secret)

		// Request id for tracing
		req.Header.Set("X-Request-Id", randomHex(16))

		// Optional: override timeout via env
		if t := os.Getenv("ENRICH_HTTP_TIMEOUT_MS"); t != "" {
			if ms, err := strconv.Atoi(t); err == nil && ms > 0 {
				httpClient.Timeout = time.Duration(ms) * time.Millisecond
			}
		}

		// Light retry on 5xx/timeouts since endpoint is idempotent (unique domain)
		retries := 2
		if v := os.Getenv("ENRICH_RETRIES"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				retries = n
			}
		}

		for attempt := 0; attempt <= retries; attempt++ {
			resp, err := httpClient.Do(req)
			if err != nil {
				if attempt == retries {
					fmt.Println("Desert enrichment request failed:", err)
					return
				}
				time.Sleep(time.Duration(150*(attempt+1)) * time.Millisecond)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				fmt.Println("Successfully triggered enrichment on Desert")
				return
			}
			// 4xx: don’t retry; 5xx: retry
			if resp.StatusCode >= 500 && attempt < retries {
				time.Sleep(time.Duration(150*(attempt+1)) * time.Millisecond)
				continue
			}
			fmt.Println("Desert enrichment API returned:", resp.StatusCode)
			return
		}
	}(domain, domainType)
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
