package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"booking-platform/internal/database"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Module struct {
	service *AuthService
	store   *storeManager
	handler *AuthHandler
}
type ModuleParams struct {
	JwtSecret           string
	DB                  *gorm.DB
	UserReaderInterface UserReaderInterface //contract it expects to read from user
	Transactor          database.TransactorInterface
}

// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "Auth Module"
}

func NewModule(params ModuleParams) *Module {
	//AuthService implements Authenticate that authenticator uses
	store := NewStoreManager([]byte(params.JwtSecret), "app_session")
	passwordRepo := NewPasswordRepository(params.DB)
	trx := database.NewGormTransactor(params.DB)
	svc := NewAuthService(params.UserReaderInterface, []byte(params.JwtSecret), store, passwordRepo, trx)
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

		r.Group(func(r chi.Router) {
			//GUESTs aren't allowed here
			r.Use(authorisation.RequiredRoles(
				core.RoleSuperAdmin,
				core.RoleAdmin,
				core.RoleCustomer,
			))

			r.Post("/password/reset", m.handler.RequestPasswordReset)
			r.Post("/password", m.handler.ResetPassword)
		})
	})
}
