package service

import (
	"encoding/json"
	"os"
	"time"

	jose "gopkg.in/square/go-jose.v2"
)

func GenerateJWT(userIDHex, orgSlug, orgType, membershipId string) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       jwtSecret,
	}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":     userIDHex,
		"org":     orgSlug,
		"orgType": orgType,
		"mid":     membershipId,
		"iat":     now,
		"exp":     now + 86400, //24 hours in seconds
	}

	// Marshal claims into JSON
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Sign the payload
	jws, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}

	// Return compact JWT string
	return jws.CompactSerialize()
}
