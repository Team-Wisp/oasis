package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const domainMapFile = "domain_map.json"

func MapDomainToOrg(domain string) string {
	// Load existing mapping
	var mapping map[string]string = make(map[string]string)

	var data []byte
	data, err := os.ReadFile(domainMapFile)
	if err == nil {
		json.Unmarshal(data, &mapping)
	}

	// Check if domain is already present in file
	var org string
	if org, exists := mapping[domain]; exists {
		return org
	}

	// if not present then fetch new mapping
	org = fetchOrgNameFromOpenAI(domain)

	// Update mapping and write back to file
	mapping[domain] = org
	updatedData, _ := json.MarshalIndent(mapping, "", "  ")
	os.WriteFile(domainMapFile, updatedData, 0644)

	return org
}

func fetchOrgNameFromOpenAI(domain string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set.")
		return "Unknown"
	}

	prompt := fmt.Sprintf("Return only the full legal organization name associated with the domain %s. No description or explanation.", domain)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":       "gpt-4o",
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  20,
		"temperature": 0.2,
	})

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return "Unknown"
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
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		orgName := strings.TrimSpace(result.Choices[0].Message.Content)
		return orgName
	}

	return "Unknown"
}
