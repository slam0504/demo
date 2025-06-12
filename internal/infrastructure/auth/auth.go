package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/google/uuid"
)

// User represents a simple user account.
type User struct {
	ID       uuid.UUID
	Username string
	Password string // stored as SHA256 hex
}

// Service manages users and sessions in memory.
type Service struct {
	mu       sync.RWMutex
	users    map[string]*User
	sessions map[string]uuid.UUID
}

// NewService creates a new auth service with a default user.
func NewService() *Service {
	s := &Service{users: make(map[string]*User), sessions: make(map[string]uuid.UUID)}
	// default user: user/password
	s.users["user"] = &User{ID: uuid.New(), Username: "user", Password: hash("password")}
	return s
}

func hash(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

// Login validates credentials and returns a session token.
func (s *Service) Login(username, password string) (string, bool) {
	s.mu.RLock()
	u, ok := s.users[username]
	s.mu.RUnlock()
	if !ok || u.Password != hash(password) {
		return "", false
	}
	token := uuid.New().String()
	s.mu.Lock()
	s.sessions[token] = u.ID
	s.mu.Unlock()
	return token, true
}

// Authenticate returns the user ID associated with a token.
func (s *Service) Authenticate(token string) (uuid.UUID, bool) {
	s.mu.RLock()
	id, ok := s.sessions[token]
	s.mu.RUnlock()
	return id, ok
}
