package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		payload := map[string]string{
			"domain":     domain,
			"domainType": domainType,
		}

		requestBody, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Failed to marshal enrichment payload:", err)
			return
		}

		url := os.Getenv("DESERT_ENRICH_URI") // internal Desert service URL
		if url == "" {
			fmt.Println("DESERT_URL is not set in env variables.")
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to create HTTP request to Desert:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to send request to Desert:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Desert enrichment API failed with status:", resp.StatusCode)
		} else {
			fmt.Println("Successfully triggered enrichment on Desert!")
		}
	}(domain, domainType)
}
