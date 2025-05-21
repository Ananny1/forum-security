🛡️ Secure Forum - GoLang
A secure, modern web forum built in Go that implements best practices in HTTPS, password hashing, session management, and rate limiting. This project serves as a learning platform for secure web development using Go.

🔐 Features
HTTPS: Secure communication using TLS (SSL certificates via server.crt and server.key).

Password Encryption: Client passwords encrypted using bcrypt.

Session Management:

Secure, HTTP-only cookies.

Server-stored sessions with UUID identifiers.

Rate Limiting:

Protects against brute-force and DoS attacks.

Goroutine-based limiter implementation.

Error Handling: Custom HTTP error responses and panic recovery.

Security Best Practices: Protection against common web vulnerabilities.

🧱 Project Structure
bash
Copy
Edit
.
├── Auth/              # Google OAuth and custom auth logic
├── database/          # DB connection and migration logic
├── handlers/          # HTTP route handlers
├── static/            # Frontend static assets (CSS, JS)
├── structs/           # Shared structs
├── Templates/         # HTML templates
├── forum.db           # SQLite3 database
├── main.go            # Main application entry
├── server.crt/key     # TLS certificate & private key
└── dockerfile         # Docker container config
🛠️ Technologies Used
Go Standard Library

SQLite3 (github.com/mattn/go-sqlite3)

bcrypt (golang.org/x/crypto/bcrypt)

UUID (github.com/google/uuid)

Gorilla Mux for routing

Docker for containerization (optional)

🚀 Getting Started
1. Clone the Repo
bash
Copy
Edit
git clone https://github.com/youruser/secure-forum.git
cd secure-forum
2. Generate TLS Certificates
bash
Copy
Edit
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
3. Build & Run
bash
Copy
Edit
go mod tidy
go run main.go
4. Visit
Open https://localhost:8080 in your browser.

⚙️ Security Details
HTTPS / TLS
Uses crypto/tls and net/http with ListenAndServeTLS.

Ensure server.crt and server.key are kept secure.

Password Encryption
go
Copy
Edit
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
Secure Sessions
Session IDs are generated using UUID.

Stored server-side in a secure session store.

Cookies:

HttpOnly

Secure

SameSite=Strict

Rate Limiting
Implemented using time.Ticker and map[IP]requests to throttle malicious users.

✅ Unit Testing
bash
Copy
Edit
go test ./...
Custom test files located in handlers/ and database/.

🧪 Example Endpoint
Create User
POST /signup

json
Copy
Edit
{
  "username": "sayed",
  "email": "sayed@example.com",
  "password": "securepassword"
}
📦 Deployment (Docker)
Dockerfile
Copy
Edit
# dockerfile

FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o forum
CMD ["./forum"]
📚 License
This project is for educational purposes only.