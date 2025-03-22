# oasis
Open-source Address &amp; Sender Identification Service

## Purpose:
Handles email domain validation, company mapping, OTP flow, and initial credential hashing

## How to run:
- in the root directory run go run cmd/server/main.go
- test /verify route
sample use:
- curl: 
curl -X POST http://localhost:8080/verify \
  -H "Content-Type: application/json" \
  -d '{"email":"user@somefreshdomain.ai"}'

- postman:
use post method for /verify with valid request

