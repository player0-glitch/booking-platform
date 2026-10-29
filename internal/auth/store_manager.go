package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

var ErrTokenNotFound = errors.New("Session Token Not Found")

const JWTTokenKey = "jwt_token"

type storeManager struct {
	session sessions.Store
	//This
	// store SessionManager
	name string
}

func NewStoreManager(session sessions.Store, name string) *storeManager {
	return &storeManager{
		session: session,
		name:    name,
	}
}

func (sm *storeManager) GetToken(r *http.Request) (string, error) {
	session, err := sm.session.Get(r, sm.name)
	if err != nil {
		return "", err
	}
	token, ok := session.Values[JWTTokenKey].(string)
	if !ok || token == "" {
		return "", ErrTokenNotFound
	}
	return token, nil
}

func (sm *storeManager) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.session.Get(r, sm.name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return sm.session.Save(r, w, session)
}

func (sm *storeManager) SaveToken(w http.ResponseWriter, r *http.Request,
	token string) error {
	session, err := sm.session.Get(r, sm.name)
	if err != nil {
		return err
	}
	session.Values[JWTTokenKey] = token
	return sm.session.Save(r, w, session)
}

// these are my older Get and Set implementations
func (sm *storeManager) Get(w http.ResponseWriter, r *http.Request, key string) (any, error) {
	store, err := sm.session.Get(r, sm.name)
	if err != nil {
		return nil, err
	}
	val, ok := store.Values[key]
	if !ok {
		return nil, errors.New("Key Not Found In Store Session")
	}
	return val, nil
}

func (sm *storeManager) Set(w http.ResponseWriter, r *http.Request, key, val string) error {
	session, err := sm.session.Get(r, sm.name)
	if err != nil {
		return errors.New("Failed To Save The Key In The Store")
	}
	session.Values[key] = val
	return sm.session.Save(r, w, session)
}

func (s *storeManager) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := s.session.Get(r, s.name)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	return s.session.Save(r, w, session)
}
