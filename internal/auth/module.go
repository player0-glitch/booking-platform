package auth

import (
	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Module struct {
	service *AuthService
	store   *StoreManager
	handler *AuthHandler
}
type ModuleParams struct {
	JwtSecret  string
	DB         *gorm.DB
	UserReader UserReader //contract it expects to read from user
}

// Contracts/Interfaces owned by the consumers
// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "auth module"
}

func NewModule(params ModuleParams) *Module {
	//AuthService implements Authenticate that authenticator uses
	store := NewStoreManager([]byte(params.JwtSecret), "app_session")
	svc := NewAuthService(params.UserReader, []byte(params.JwtSecret), store)
	return &Module{
		service: svc,
		store:   store,
		handler: NewAuthHandler(svc),
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
