package sessions

import (
	"booking-platform/internal/database"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	gorillaSessions "github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type SessionStore struct {
	db          *gorm.DB
	codecs      []securecookie.Codec
	transaction database.TransactorInterface
	//These are publically accessible
	Options *gorillaSessions.Options
}

func NewSesssionStore(
	db *gorm.DB,
	hashKey []byte,
	blockKey []byte,
	trx database.TransactorInterface,
) *SessionStore {
	return &SessionStore{
		db: db,
		codecs: []securecookie.Codec{
			securecookie.New(hashKey, blockKey),
		},
		Options: &gorillaSessions.Options{
			Path:     "/",
			MaxAge:   3600 * 2, /*2hrs*/
			Secure:   false,    /*set this to true in prod for HTTPS*/
			SameSite: http.SameSiteLaxMode,
		},
		transaction: trx,
	}
}
func (s *SessionStore) New(
	r *http.Request,
	name string,
) (*gorillaSessions.Session, error) {
	session := gorillaSessions.NewSession(s, name)
	session.Options = s.Options
	session.IsNew = true
	session.Values = make(map[any]any)
	return session, nil
}
func (s *SessionStore) Get(
	r *http.Request,
	name string,
) (*gorillaSessions.Session, error) {
	return s.load(r, name)
}

func (s *SessionStore) Save(
	r *http.Request,
	w http.ResponseWriter,
	session *gorillaSessions.Session,
) error {
	return nil
}
func (s *SessionStore) load(
	r *http.Request,
	name string,
) (*gorillaSessions.Session, error) {
	session, err := s.New(r, name)
	if err != nil {
		return nil, err
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}
	var id string
	if err := securecookie.DecodeMulti(
		name, cookie.Value, &id, s.codecs...,
	); err != nil {
		return session, nil
	}
	sessionId, err := uuid.Parse(id)
	if err != nil {
		return session, nil
	}
	var record AppSession
	err = s.db.WithContext(r.Context()).
		Where("id = ? AND expires_at > CURRENT_TIMESTAMP", &sessionId).
		Take(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, nil
		}
		return nil, err
	}
	var vals map[string]any
	if err := json.Unmarshal(record.Data, &vals); err != nil {
		return nil, err
	}
	session.Values = make(map[any]any, len(vals))
	for k, v := range vals {
		session.Values[k] = v
	}

	session.ID = id
	session.IsNew = false
	session.Options = s.Options
	return session, nil
}
