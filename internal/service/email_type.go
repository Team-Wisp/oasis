package service

import "strings"

func GetDomainType(domain string) string {
	publicEmailDomains := map[string]bool{
		// "gmail.com":      true,
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
