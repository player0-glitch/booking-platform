package auth

import (
	"strings"
	"sync"
	"time"
)

// Json Requests This packages requires/uses

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
		info = &loginAttempInfo{
			lockoutEnd:    time.Now(),
			attempts:      1,
			requiresReset: false,
		}
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
