package auth

import (
	chi "github.com/go-chi/chi/v5"
)

type Module struct {
	// session sessions.Session
	service *AuthService
	store   *StoreManager
	handler *AuthHandler
}

// Contracts/Interfaces owned by the consumers
// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "auth module"
}

func NewModule(jwtSecret string) *Module {
	jwtBytes := []byte(jwtSecret)
	//AuthService implements Authenticate that authenticator uses
	srv := NewAuthService(jwtBytes)
	store := NewStoreManager(jwtBytes, "app_session")
	return &Module{
		// session: s,
		service: srv,
		store:   store,
		handler: NewAuthHandler(srv, store),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		//These endpoints should  be standalone
		//They both are concerned about auth but carry totaly different responsibilities
		r.Post("/login", m.handler.Login)
		// r.Post("/logout", m.handler.Logout)
	})
}
