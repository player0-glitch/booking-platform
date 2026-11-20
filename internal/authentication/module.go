package authentication

import (
	"booking-platform/internal/database"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// var (
// 	hashKeySecret  = os.Getenv("SESSION_KEY_SECRET")
// 	blockKeySecret = os.Getenv("SESSION_BLOCK_SECRET")
// 	jwtSecret      = os.Getenv("JWT_SECRET")
// )

type Module struct {
	service    *Service
	handler    *Handler
	Middleware *AuthenticationMiddleware
}
type ModuleParams struct {
	DB                  *gorm.DB
	UserReaderInterface UserReaderInterface
	HandlerStoreManager handlerSessionInterface
	JWTGenerator        JwtTokenGenerator
}

// Implement this and RegisterRoutes so it can be plugged/injected into the app
func (m *Module) Name() string {
	return "Auth Module"
}

func NewModule(params ModuleParams) *Module {
	// Transactor for db transactions
	trx := database.NewGormTransactor(params.DB)
	// Password Repo Interface
	passwordRepo := NewPasswordRepository(params.DB)
	//Service
	svc := NewService(
		params.UserReaderInterface,
		params.JWTGenerator,
		passwordRepo,
		trx)

	// Lastly Create the AppMiddleware
	return &Module{
		service: svc,
		handler: NewHandler(svc, params.HandlerStoreManager),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/authentication", func(r chi.Router) {
		// These endpoints should  be standalone
		// They both are concerned about authentication but carry totaly different responsibilities
		r.Post("/login", m.handler.Login)
		r.Post("/logout", m.handler.Logout)

		r.Group(func(protected chi.Router) {
			protected.Post("/password/reset", m.handler.RequestPasswordReset)
			protected.Post("/password", m.handler.ResetPassword)
		})
	})
}
