package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type DomainInfo struct {
	Domain    string    `bson:"domain"`
	Name      string    `bson:"name"`
	Type      string    `bson:"type"`
	CreatedAt time.Time `bson:"createdAt"`
}

func GetOrInitDomain(domain, domainType string) {
	go func(domain, domainType string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := MongoDatabase.Collection("organizations")

		// Check if domain already exists
		var existing DomainInfo
		err := collection.FindOne(ctx, map[string]interface{}{"domain": domain}).Decode(&existing)
		if err == nil {
			// Already exists, no need to fetch
			return
		}

		// If not found, fetch using OpenAI
		org := fetchOrgNameFromOpenAI(domain)
		if org == "" {
			org = "Unknown"
		}

		// Insert into MongoDB
		newDomain := DomainInfo{
			Domain:    domain,
			Name:      org,
			Type:      domainType,
			CreatedAt: time.Now(),
		}

		_, err = collection.InsertOne(ctx, newDomain)
		if err != nil {
			fmt.Println("Failed to insert new domain info:", err)
		} else {
			fmt.Printf("Cached new domain in MongoDB: %s => %s (%s)\n", domain, org, domainType)
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
