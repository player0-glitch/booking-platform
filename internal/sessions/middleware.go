package sessions

import "net/http"

type SessionMiddleware struct {
	manager *Manager
}

func NewSessionMiddleware(manager *Manager) *SessionMiddleware {
	return &SessionMiddleware{
		manager: manager,
	}
}

func (m *SessionMiddleware) EnsureAnonymous(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := m.manager.Get(r)
		if err != nil {
			http.Error(w, "unable to load session", http.StatusInternalServerError)
			return
		}
		// Authenticated User (already exists that's how they logged in to be 'authenticated' duh)
		if sess.Values[SessionType] != "anonymous" {
			sess.Values[SessionType] = "anonymous"

			if err := m.manager.Save(w, r, sess); err != nil {
				http.Error(w, "unable to save token", http.StatusInternalServerError)
				return
			}
		} else if sess.IsNew {
			if err := m.manager.Save(w, r, sess); err != nil {
				http.Error(w, "unable to save anonymous session", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
