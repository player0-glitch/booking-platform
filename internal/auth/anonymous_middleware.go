package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const SessionTypeKey = "session_type"

type AnonymousMiddelware struct {
	store sessions.Store
	name  string
}

func NewAnonymousMiddelware(store sessions.Store, name string) *AnonymousMiddelware {
	return &AnonymousMiddelware{
		store: store,
		name:  name,
	}
}
func (m *AnonymousMiddelware) Ensure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, m.name)
		if err != nil {
			http.Error(w, "Unable to load session", http.StatusInternalServerError)
			return
		}
		//Provision a guest by marking a session as anonymous
		if session.Values[SessionTypeKey] != "anonymous" {
			session.Values[SessionTypeKey] = "anonymous"

			if err := m.store.Save(r, w, session); err != nil {
				http.Error(w, "Unable to persist anonymous session",
					http.StatusInternalServerError)
				return
			}
		} else if session.IsNew {
			if err := m.store.Save(r, w, session); err != nil {
				http.Error(w, "Unable to persist anonymous session", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
