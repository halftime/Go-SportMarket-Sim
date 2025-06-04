package server

import (
	"sync"
	"time"
)

type BaseReply struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type ErrorReply struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Data   any    `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionData struct {
	Username    string `json:"username"`
	Client_type string `json:"client_type"`
}

// Server represents the application server with configuration
type Server struct {
	URL         string
	APILimit    int
	RateLimiter map[string]*RateLimit
	mu          sync.Mutex
}

// RateLimit tracks API usage for rate limiting
type RateLimit struct {
	Count      int
	ResetTime  time.Time
	WindowSize time.Duration
}

// DefaultServer is the shared instance that can be used across different files
var DefaultServer = &Server{
	URL:         "http://localhost:8080/v1",
	APILimit:    5,
	RateLimiter: make(map[string]*RateLimit),
}

// IsRateLimited checks if a client has exceeded their API limit
func (s *Server) IsRateLimited(clientID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	limit, exists := s.RateLimiter[clientID]

	if !exists {
		s.RateLimiter[clientID] = &RateLimit{
			Count:      1,
			ResetTime:  now.Add(time.Minute * 5), // reset after 5 minute
			WindowSize: time.Hour,
		}
		return false
	}

	// Reset counter if time window has elapsed
	if now.After(limit.ResetTime) {
		limit.Count = 1
		limit.ResetTime = now.Add(limit.WindowSize)
		return false
	}

	// Increment counter and check against limit
	limit.Count++
	return limit.Count > s.APILimit
}

// GetRemainingLimit returns the number of API calls remaining for a client
func (s *Server) GetRemainingLimit(clientID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	limit, exists := s.RateLimiter[clientID]
	if !exists {
		return s.APILimit
	}

	if time.Now().After(limit.ResetTime) {
		return s.APILimit
	}

	remaining := s.APILimit - limit.Count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ResetRateLimits clears all rate limiting data
func (s *Server) ResetRateLimits() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RateLimiter = make(map[string]*RateLimit)
}

func (s *Server) AuthenticateLogin(loginreq LoginRequest) bool {
	// This is a placeholder for actual authentication logic
	// In a real application, you would check the credentials against a database
	return loginreq.Username == "testuser" && loginreq.Password == "testpass"
}

func (s *Server) GenerateClientSessionID(username string) string {
	// This is a placeholder for generating a session ID
	// In a real application, you would generate a secure session ID
	return "session_" + username
}
