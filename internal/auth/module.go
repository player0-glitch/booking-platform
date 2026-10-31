package auth

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"booking-platform/internal/database"
	"os"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var (
	hashKeySecret  = os.Getenv("SESSION_KEY_SECRET")
	blockKeySecret = os.Getenv("SESSION_BLOCK_SECRET")
	jwtSecret      = os.Getenv("JWT_SECRET")
)

type Module struct {
	service       *AuthService
	store         *storeManager
	handler       *AuthHandler
	AppMiddleware *AppMiddleware
}
type ModuleParams struct {
	DB                  *gorm.DB
	UserReaderInterface UserReaderInterface //contract it expects to read from user
	Transactor          database.TransactorInterface
}

// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "Auth Module"
}

func NewModule(params ModuleParams) *Module {
	//Transactor for db transactions
	trx := database.NewGormTransactor(params.DB)
	//Session Wrapper for Gorilla Sessions
	session := NewSessionManager(
		params.DB,
		[]byte(hashKeySecret),
		[]byte(blockKeySecret),
		trx)
	//Session Store Manager (interacts with DB)
	store := NewStoreManager(session, "app_session")
	//Instantiate JWT Service
	jwtSvc := NewJWTService(jwtSecret)
	//Password Repo Interface
	passwordRepo := NewPasswordRepository(params.DB)
	//User Service
	svc := NewAuthService(
		params.UserReaderInterface,
		jwtSvc,
		store,
		passwordRepo,
		trx)
	//Create jwt service

	//Lastly Create the AppMiddleware
	appMW := NewAppMiddleware(store, "app_session", jwtSvc)
	return &Module{
		service:       svc,
		store:         store,
		handler:       NewAuthHandler(svc),
		AppMiddleware: appMW,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		//These endpoints should  be standalone
		//They both are concerned about auth but carry totaly different responsibilities
		r.Post("/login", m.handler.Login)
		r.Post("/logout", m.handler.Logout)

		r.Group(func(protected chi.Router) {
			//GUESTs aren't allowed here
			protected.Use(authorisation.RequiredRoles(
				core.RoleSuperAdmin,
				core.RoleAdmin,
				core.RoleCustomer,
			))

			protected.Post("/password/reset", m.handler.RequestPasswordReset)
			protected.Post("/password", m.handler.ResetPassword)
		})
	})
}
