package auth

import (
	"strings"
	"sync"
	"time"
)

// should be onto core module
type ErrorResponse struct {
	Code    int64  `json:"status"`
	Message string `json:"message"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}

// godoc
// object that will be used track of account lock out
type loginAttempInfo struct {
	attempts      int
	lockoutEnd    time.Time
	requiresReset bool
}

// godoc
// Concurren Store to mimic C#'s Concurrent Dictionary
type concurrentLoginAttemptStore struct {
	//I think this is a read write mutex lock
	mu   sync.RWMutex
	data map[string]*loginAttempInfo
}

func NewConcurrentLoginAttemptStore() *concurrentLoginAttemptStore {
	return &concurrentLoginAttemptStore{
		data: make(map[string]*loginAttempInfo),
	}
}

// godoc
// Get Or Add email to login attempts concurrent store. This is used to track login
// attemps of a user
func (s *concurrentLoginAttemptStore) GetOrAdd(email string) *loginAttempInfo {
	key := strings.ToLower(email)
	s.mu.Lock()
	defer s.mu.Unlock()

	info, exists := s.data[key]
	if !exists {
		//info don't exists, add to the store
		info = &loginAttempInfo{}
		s.data[key] = info
	}

	return info
}

// godoc
// Remove the email that was attempted for loging in if reset is not required
func (s *concurrentLoginAttemptStore) Remove(email string) {
	key := strings.ToLower(email)

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
