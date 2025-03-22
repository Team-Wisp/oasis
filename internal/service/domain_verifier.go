package service

import (
	"log"
	"net"
)

// Lookup to check if the domain is config for receiving mail
func CheckMX(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("MX lookup failed for domain: %s, error: %v", domain, err)
		return false
	}
	log.Printf("MX records found for domain: %s", domain)
	return len(mxRecords) > 0
}
