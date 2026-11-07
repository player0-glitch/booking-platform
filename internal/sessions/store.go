// Package sessions acts as a wrapper for gorilla/sessions to persist on Postgres
package sessions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"booking-platform/internal/core"
	"booking-platform/internal/database"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	gorillaSessions "github.com/gorilla/sessions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionStore struct {
	db          *gorm.DB
	codecs      []securecookie.Codec
	transaction database.TransactorInterface
	// These are publically accessible
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
	// cannot override an existing session when saving
	if session.Options.MaxAge < 0 {
		if session.ID != "" {
			id, err := uuid.Parse(session.ID)
			if err != nil {
				return nil
			}

			if err := s.db.WithContext(r.Context()).Unscoped().
				Delete(&AppSession{}, "id = ? ", id).Error; err != nil {
				return err
			}
		}
		// Set a new cookie after storing the session ID in the db
		http.SetCookie(w, gorillaSessions.NewCookie(session.Name(), "",
			session.Options))
		return nil
	}
	values := make(map[string]any, len(session.Values))

	for k, v := range session.Values {
		strKey, ok := k.(string)
		if !ok {
			return fmt.Errorf("session key must be a string")
		}
		values[strKey] = v
	}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to encode session data %w ", err)
	}
	id, err := uuid.Parse(session.ID)
	if err != nil {
		id = uuid.New()
		session.ID = id.String()
	}
	record := AppSession{
		base:      core.Base{ID: id},
		Data:      data,
		ExpiresAt: time.Now().Add(time.Duration(session.Options.MaxAge) * time.Second),
	}
	err = s.db.WithContext(r.Context()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"date":       record.Data,
			"expires_at": record.ExpiresAt,
			"updated_at": record.base.UpdatedAt,
		}),
	}).Create(&record).Error
	if err != nil {
		return fmt.Errorf("persisting state failed: %v", err)
	}

	encoded, err := securecookie.EncodeMulti(
		session.Name(),
		session.ID,
		s.codecs...,
	)
	if err != nil {
		return fmt.Errorf("failed encoding a session cookie: %w", err)
	}
	http.SetCookie(w, gorillaSessions.NewCookie(
		session.Name(),
		encoded,
		session.Options,
	))
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
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return session, nil
	}
	var record AppSession
	err = s.db.WithContext(r.Context()).
		Where("id = ? AND expires_at > CURRENT_TIMESTAMP", &sessionID).
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
