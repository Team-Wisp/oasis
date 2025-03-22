package service

import (
	"log"
	"net"
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
