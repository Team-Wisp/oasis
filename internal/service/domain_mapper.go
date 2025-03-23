package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type DomainInfo struct {
	Organization string `json:"organization"`
	Type         string `json:"type"`
}

const domainMapFile = "domain_map.json"

var domainMapLock sync.Mutex

func GetOrInitDomain(domain, domainType string) {
	// Lock for safe read/write to domain_map.json
	domainMapLock.Lock()
	defer domainMapLock.Unlock()

	// Load existing data
	mapping := make(map[string]DomainInfo)
	data, err := os.ReadFile(domainMapFile)
	if err == nil {
		_ = json.Unmarshal(data, &mapping)
	}

	// If already exists, do nothing
	if _, exists := mapping[domain]; exists {
		return
	}

	// Background goroutine to fetch & write mapping
	go func(domain, domainType string) {
		org := fetchOrgNameFromOpenAI(domain)
		if org == "" {
			org = "Unknown"
		}

		domainMapLock.Lock()
		defer domainMapLock.Unlock()

		// Reload to avoid overwriting concurrent updates
		updatedMapping := make(map[string]DomainInfo)
		data, err := os.ReadFile(domainMapFile)
		if err == nil {
			_ = json.Unmarshal(data, &updatedMapping)
		}

		updatedMapping[domain] = DomainInfo{
			Organization: org,
			Type:         domainType,
		}

		// Write back to file
		updatedData, _ := json.MarshalIndent(updatedMapping, "", "  ")
		_ = os.WriteFile(domainMapFile, updatedData, 0644)

		fmt.Printf("✅ Cached new domain: %s => %s (%s)\n", domain, org, domainType)
	}(domain, domainType)
}

func fetchOrgNameFromOpenAI(domain string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  OPENAI_API_KEY not set.")
		return ""
	}

	prompt := fmt.Sprintf("Return only the full legal organization name associated with the domain %s. No description or explanation.", domain)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":       "gpt-4o",
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  20,
		"temperature": 0.2,
	})

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // this is like await => wait until this call is done
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		return strings.TrimSpace(result.Choices[0].Message.Content)
	}

	return ""
}
