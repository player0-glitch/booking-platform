package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"context"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrSessionNotFound = errors.New("Session Not Found")
	ErrSessionExpired  = errors.New("Session Found To Be Expired")
	ErrKeyNotFound     = errors.New("Key Not Found In Sessions")
	ErrInvalidSession  = errors.New("Invalide Session")
)

const (
	jwtKey       = "jwt_token"
	anonymousKey = "anonymous_token"
)

type SessionManager struct {
	db         *gorm.DB
	Options    *sessions.Options
	Codecs     []securecookie.Codec
	transactor database.TransactorInterface
}

func NewSessionManager(db *gorm.DB, hashKey []byte, blockKey []byte, trx database.TransactorInterface) *SessionManager {
	return &SessionManager{
		db: db,
		Codecs: []securecookie.Codec{
			securecookie.New(hashKey, blockKey),
		},
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   3600 * 2, /*lives for 2hrs*/
			HttpOnly: true,
			Secure:   false, /*set this to true in prod for HTTPS*/
			SameSite: http.SameSiteLaxMode,
		},
		transactor: trx,
	}
}

func (s *SessionManager) New(
	r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = s.Options
	session.IsNew = true
	session.Values = make(map[any]any)
	return session, nil
}

func (s *SessionManager) Get(
	r *http.Request, name string) (*sessions.Session, error) {
	// return sessions.GetRegistry(r).Get(s, name)
	return s.load(r, name)
}
func (s *SessionManager) Save(
	r *http.Request, w http.ResponseWriter, session *sessions.Session) error {

	//make sure our store exists
	if session.Options == nil {
		session.Options = s.Options
	}

	if session.Options.MaxAge < 0 {
		//expired sesssion
		if session.ID != "" {
			id, err := uuid.Parse(session.ID)
			if err != nil {
				return err
			}
			//Delete the expired session
			if err := s.db.WithContext(r.Context()).
				Delete(&AppSession{}, "id = ?", id).
				Error; err != nil {
				return err
			}
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(),
			"",
			session.Options,
		))
		return nil
	}

	id := session.ID
	if id == "" {
		id = uuid.NewString()
	}

	sessionId, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidSession
	}
	vals := make(map[string]any, len(session.Values))
	for k, v := range session.Values {
		stringKey, ok := k.(string)
		if !ok {
			return errors.New("Session Key Must Be A String")
		}
		vals[stringKey] = v
	}

	data, err := json.Marshal(vals)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(session.Options.MaxAge) * time.Second)

	record := AppSession{
		Base: core.Base{ID: sessionId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		Data:      data,
		ExpiresAt: expiresAt,
	}
	// Keep an existing authenticated user association on session updates.
	// New login/session-association logic may set this field explicitly.
	var txErr = s.transactor.Transaction(r.Context(),
		func(txCtx context.Context) error {
			var existingRecord AppSession
			dbTrx := database.DbFromTransaction(r.Context(), s.db)
			findErr := dbTrx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", sessionId).
				Take(&existingRecord).Error

			if findErr == nil {
				record.Base.CreatedAt = existingRecord.ExpiresAt
				record.ClaimedByUserId = existingRecord.ClaimedByUserId
				return dbTrx.Model(&existingRecord).
					Updates(map[string]any{
						"data":       record.Data,
						"expires_at": record.ExpiresAt,
						"updated_at": record.Base.UpdatedAt,
					}).Error
			}
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return findErr
			}
			//automatically commit when done
			return dbTrx.Create(&record).Error
		})
	if txErr != nil {
		return txErr
	}

	encoded, err := securecookie.EncodeMulti(
		session.Name(),
		sessionId.String(),
		s.Codecs...,
	)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(
		session.Name(),
		encoded,
		session.Options,
	))

	session.ID = sessionId.String()
	session.IsNew = false

	return nil
}
func (s *SessionManager) load(r *http.Request, name string) (*sessions.Session, error) {
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
		name, cookie.Value, &id, s.Codecs...,
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
