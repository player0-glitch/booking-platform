package core

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type StoreManager struct {
	store      sessions.Store
	cookieName string
}

// NewStoreManager go doc
// This will be used as a session manager and passed around the codebase
func NewStoreManager(store sessions.Store, coockieName string) *StoreManager {
	return &StoreManager{
		store:      store,
		cookieName: coockieName,
	}
}

func (s *StoreManager) Get(r *http.Request, key string) (any, error) {
	store, err := s.store.Get(r, s.cookieName)

	if err != nil {
		return nil, err
	}
	val, ok := store.Values[key]
	if !ok {
		return nil, errors.New("Key not found in store session")
	}
	return val, nil
}

func (s *StoreManager) Set(w http.ResponseWriter, r *http.Request, key string, val any) error {
	store, _ := s.store.Get(r, s.cookieName)
	store.Values[key] = val
	return store.Save(r, w)
}
func (s *StoreManager) Clear(w http.ResponseWriter, r *http.Request) error {
	store, _ := s.store.Get(r, s.cookieName)
	//invalidate the current session from the store and make it expire
	store.Options.MaxAge = -1
	return store.Save(r, w)
}

func (s *StoreManager) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	store, _ := s.store.Get(r, s.cookieName)
	delete(store.Values, key)
	return store.Save(r, w)
}
