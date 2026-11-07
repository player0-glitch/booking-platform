package sessions

import (
	"errors"
	"net/http"

	gorillaSessions "github.com/gorilla/sessions"
)

var (
	ErrSessionNotFound       = errors.New("session not found")
	ErrKeyNotFound           = errors.New("key not found in sessions")
	ErrJwtKeyNotFound  error = errors.New("jwt token not found in session")
)

const (
	// remember why we're here
	SessionName = "booking_session"
	JWTToken    = "jwt_token"
	SessionType = "session_type"
)

type Manager struct {
	store gorillaSessions.Store
	name  string
}

func NewManager(store gorillaSessions.Store, name string) *Manager {
	return &Manager{
		store: store,
		name:  name,
	}
}

func (m *Manager) Get(r *http.Request) (*gorillaSessions.Session, error) {
	return m.store.Get(r, m.name)
}

func (m *Manager) GetValue(r *http.Request, key string) (any, error) {
	// get the session in this context
	session, err := m.Get(r)
	if err != nil {
		return nil, ErrKeyNotFound
	}
	val, exists := session.Values[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	return val, nil
}

func (m *Manager) Set(
	w http.ResponseWriter,
	r *http.Request,
	key string,
	val any,
) error {
	// get the session in this context
	session, err := m.Get(r)
	if err != nil {
		return ErrSessionNotFound
	}
	session.Values[key] = val

	return m.store.Save(r, w, session)
}

func (m *Manager) Delete(
	w http.ResponseWriter,
	r *http.Request,
	key string,
) error {
	session, err := m.Get(r)
	if err != nil {
		return ErrSessionNotFound
	}
	delete(session.Values, key)

	return m.store.Save(r, w, session)
}

// GetToken Gets the JWT Token from session
func (m *Manager) GetToken(r *http.Request) (string, error) {
	val, err := m.GetValue(r, JWTToken)
	if err != nil {
		return "", ErrKeyNotFound
	}
	token, ok := val.(string)
	if !ok {
		return "", ErrJwtKeyNotFound
	}
	return token, nil
}

func (m *Manager) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := m.Get(r)
	if err != nil {
		return err
	}

	// invalidate this session
	session.Options.MaxAge = -1

	return m.store.Save(r, w, session)
}

func (m *Manager) Save(
	w http.ResponseWriter,
	r *http.Request,
	session *gorillaSessions.Session,
) error {
	return m.store.Save(r, w, session)
}

func (m *Manager) Rotate(
	w http.ResponseWriter,
	r *http.Request,
) (*gorillaSessions.Session, error) {
	oldSess, err := m.Get(r)
	if err != nil {
		return nil, err
	}

	//Preserve the existing sessino values
	values := make(map[any]any, len(oldSess.Values))
	for k, v := range oldSess.Values {
		values[k] = v
	}

	//invalidate the old session (move this operation to DB at some point)
	oldSess.Options.MaxAge = -1

	if err := m.store.Save(r, w, oldSess); err != nil {
		return nil, err
	}

	//immediately create a new session
	newSess, err := m.store.New(r, m.name)
	if err != nil {
		return nil, err
	}

	newSess.Values = oldSess.Values
	newSess.IsNew = true

	return newSess, nil
}
