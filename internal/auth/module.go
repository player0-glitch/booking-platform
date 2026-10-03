package auth

import (
	chi "github.com/go-chi/chi/v5"
)

type Module struct {
	service *AuthService
	store   *StoreManager
	handler *AuthHandler
}
type ModuleParams struct {
	JwtSecret   string
	AuthService *AuthService
}

// Contracts/Interfaces owned by the consumers
// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "auth module"
}

func NewModule(params ModuleParams) *Module {
	//AuthService implements Authenticate that authenticator uses
	store := NewStoreManager([]byte(params.JwtSecret), "app_session")
	return &Module{
		service: params.AuthService,
		store:   store,
		handler: NewAuthHandler(params.AuthService, store),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		//These endpoints should  be standalone
		//They both are concerned about auth but carry totaly different responsibilities
		r.Post("/login", m.handler.Login)
		r.Post("/logout", m.handler.Logout)
	})
}
