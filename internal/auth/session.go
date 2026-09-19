// Session Manager Wrapper
package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type StoreManager struct {
	store sessions.Store
	name  string
}

func NewStoreManager(storeSecret []byte, name string) *StoreManager {
	cookieStore := sessions.NewCookieStore(storeSecret)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 2, /*lives for 2hrs*/
		HttpOnly: true,
		Secure:   false, /*set this to true in prod for HTTPS*/
		SameSite: http.SameSiteLaxMode,
	}

	return &StoreManager{
		store: cookieStore,
		name:  name,
	}
}

func (sm *StoreManager) GetToken(r *http.Request) (string, error) {
	session, err := sm.store.Get(r, sm.name)
	if err != nil {
		return "", err
	}
	token, ok := session.Values["jwt_token"].(string)
	if !ok || token == "" {
		return "", errors.New("Session Token Not Found")
	}
	return token, nil
}

func (sm *StoreManager) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.store.Get(r, sm.name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func (sm *StoreManager) SaveToken(w http.ResponseWriter, r *http.Request,
	token string) error {
	session, err := sm.store.Get(r, sm.name)
	if err != nil {
		return err
	}
	session.Values["jwt_token"] = token

	return session.Save(r, w)
}

// these are my older Get and Set implementations
func (sm *StoreManager) Get(w http.ResponseWriter, r *http.Request, key string) (any, error) {
	store, err := sm.store.Get(r, sm.name)

	val, ok := store.Values[key]
	if !ok {
		return nil, errors.New("Key Not Found In Store Session")
	}
	return val, err
}

func (sm *StoreManager) Set(w http.ResponseWriter, r *http.Request, key, val string) error {
	store, err := sm.store.Get(r, sm.name)
	if err != nil {
		return errors.New("Failed To Save The Key In The Store")
	}
	store.Values[key] = val
	return store.Save(r, w)
}

func (s *StoreManager) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	store, _ := s.store.Get(r, s.name)
	delete(store.Values, key)
	return store.Save(r, w)
}
