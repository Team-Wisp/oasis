package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type DomainInfo struct {
	Domain    string    `bson:"domain"`
	Org_Name  string    `bson:"org_name"`
	Org_Type  string    `bson:"org_type"`
	CreatedAt time.Time `bson:"createdAt"`
}

func GetOrInitDomain(domain, domainType string) {
	go func(domain, domainType string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := MongoDatabase.Collection("organizations")

		// Sanitize domain input
		var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)

		if !domainRegex.MatchString(domain) {
			sanitizedDomain := strings.ReplaceAll(domain, "\n", "")
			sanitizedDomain = strings.ReplaceAll(sanitizedDomain, "\r", "")
			fmt.Println("Invalid domain format:", sanitizedDomain)
			return
		}

		// Check if domain already exists
		var existing DomainInfo
		err := collection.FindOne(ctx, map[string]interface{}{"domain": domain}).Decode(&existing)
		if err == nil {
			// Already exists, no need to fetch
			return
		}

		// If not found, fetch using OpenAI
		org_name := fetchOrgNameFromOpenAI(domain)
		if org_name == "" {
			org_name = "Unknown"
		}

		// Insert into MongoDB
		newDomain := DomainInfo{
			Domain:    domain,
			Org_Name:  org_name,
			Org_Type:  domainType,
			CreatedAt: time.Now(),
		}

		_, err = collection.InsertOne(ctx, newDomain)
		if err != nil {
			fmt.Println("Failed to insert new domain info:", err)
		} else {
			sanitizedDomain := strings.ReplaceAll(domain, "\n", "")
			sanitizedDomain = strings.ReplaceAll(sanitizedDomain, "\r", "")
			fmt.Printf("Cached new domain in MongoDB: %s => %s (%s)\n", sanitizedDomain, org_name, domainType)
		}
	}(domain, domainType)
}

func fetchOrgNameFromOpenAI(domain string) string {
	fmt.Println("Making OpenAI Request for new Domain Mapping")
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

	// Defines a Go struct that matches the structure of the JSON response from OpenAI's API.
	var result struct {
		Choices []struct { // Choices is an array of objects that contain the message. In this case we only care about the first one. As default it is 1 but can be more.
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	fmt.Println("OpenAI response status:", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Failed to parse OpenAI response:", err)
		return ""
	}

	if len(result.Choices) > 0 {
		return strings.TrimSpace(result.Choices[0].Message.Content)
	}

	return ""
}
