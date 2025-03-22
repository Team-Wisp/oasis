package service

import (
	"log"
	"net"
	"strings"
)

// Lookup to check if the domain is config for receiving mail
func CheckMX(domain string) bool {
	mxRecords, err := net.LookupMX(domain) // This mxRecrods will be a list of pointers to the mailservers for the domain
	if err != nil {
		log.Printf("MX lookup failed for domain: %s, error: %v", domain, err)
		return false
	}
	log.Printf("MX records found for domain: %s", domain)
	return len(mxRecords) > 0
}

func GetDomainType(domain string) string {
	publicEmailDomains := map[string]bool{
		"gmail.com":      true,
		"yahoo.com":      true,
		"outlook.com":    true,
		"hotmail.com":    true,
		"icloud.com":     true,
		"protonmail.com": true,
		"gmx.com":        true,
		"mail.com":       true,
	}

	domain = strings.ToLower(domain)

	// Check public emails first
	if publicEmailDomains[domain] {
		return "public"
	}

	// Check for educational domains (multi-suffix)
	eduSuffixes := []string{
		".edu", ".edu.", ".ac.", ".k12.", ".school", ".college", ".university",
	}

	for _, suffix := range eduSuffixes {
		if strings.Contains(domain, suffix) {
			return "educational"
		}
	}

	return "corporate"
}
