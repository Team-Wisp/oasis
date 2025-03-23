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

## RUN using Docker 
Clone the repo and run 
`docker build -t oasis-app .  `
once the image is created. 
```
create a .env to setup these ENV vars 
OPENAI_API_KEY
EMAIL_SENDER
EMAIL_PASSWORD
SMTP_HOST
SMTP_PORT
REDIS_URL
```

and then run it !!!
`docker run --rm -p 8080:8080 --env-file .env oasis-app`


