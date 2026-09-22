package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Auth handles OTP-based authentication and bearer token management.
type Auth struct {
	secret string
	mu     sync.RWMutex
	otps   map[string]otpEntry
	tokens map[string]tokenEntry
}

type otpEntry struct {
	email   string
	code    string
	expires time.Time
}

type tokenEntry struct {
	email     string
	createdAt time.Time
}

// New creates a new Auth instance.
func New(secret string) *Auth {
	return &Auth{
		secret: secret,
		otps:   make(map[string]otpEntry),
		tokens: make(map[string]tokenEntry),
	}
}

// RequestOTP generates a 6-digit OTP for the given email.
// In dev mode (no SMTP), the OTP is returned for logging.
func (a *Auth) RequestOTP(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is required")
	}

	code := generateOTP()
	a.mu.Lock()
	a.otps[email] = otpEntry{
		email:   email,
		code:    code,
		expires: time.Now().Add(5 * time.Minute),
	}
	a.mu.Unlock()

	// In production, send email here. For dev, return the code.
	return code, nil
}

// VerifyOTP validates the OTP and issues a bearer token.
func (a *Auth) VerifyOTP(email, code string) (string, error) {
	a.mu.RLock()
	entry, ok := a.otps[email]
	a.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("no OTP requested for this email | hint: call POST /auth/request with email first")
	}

	if time.Now().After(entry.expires) {
		a.mu.Lock()
		delete(a.otps, email)
		a.mu.Unlock()
		return "", fmt.Errorf("OTP expired | hint: request a new OTP via POST /auth/request")
	}

	if entry.code != code {
		return "", fmt.Errorf("invalid OTP code | hint: check the 6-digit code and try again")
	}

	// Generate bearer token
	token := generateToken(a.secret, email)

	a.mu.Lock()
	delete(a.otps, email)
	a.tokens[token] = tokenEntry{
		email:     email,
		createdAt: time.Now(),
	}
	a.mu.Unlock()

	return token, nil
}

// ValidateToken checks if a bearer token is valid.
func (a *Auth) ValidateToken(token string) bool {
	a.mu.RLock()
	_, ok := a.tokens[token]
	a.mu.RUnlock()
	return ok
}

// Middleware checks for a valid bearer token.
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		path := r.URL.Path
		if path == "/help" || path == "/.well-known/agent.md" || path == "/health" || path == "/mcp" {
			next.ServeHTTP(w, r)
			return
		}

		// Skip auth for public random endpoints (GET only)
		if r.Method == http.MethodGet && (path == "/" || strings.HasPrefix(path, "/random/") || strings.HasPrefix(path, "/password") || strings.HasPrefix(path, "/lorem") || strings.HasPrefix(path, "/uuid") || strings.HasPrefix(path, "/color") || strings.HasPrefix(path, "/coin") || strings.HasPrefix(path, "/dice")) {
			next.ServeHTTP(w, r)
			return
		}

		// Check auth header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "error: missing auth token | hint: call POST /auth/request with email to get an OTP, then POST /auth/verify to get a bearer token\n")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !a.ValidateToken(token) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "error: invalid or expired token | hint: request a new OTP via POST /auth/request and verify via POST /auth/verify\n")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateOTP() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	code := fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000)
	return code
}

func generateToken(secret, email string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(email))
	h.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(h.Sum(nil))
}
