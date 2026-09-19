package user

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler

// Module Structure it self
type Module struct {
	userService *UserService
	// Controllers Controllers
	userHandler    *UserController
	authMiddleware Middleware
}

// Implement the applications's Module interface
func (m *Module) Name() string {
	return "user module"
}

// Any external dependencies (like from the core app) needed by
// this module
type ModuleParams struct {
	DB             *gorm.DB
	AuthMiddleware Middleware
}

func NewModule(params ModuleParams) *Module {
	//register repos
	userRepo := NewUserRepo(params.DB)
	//register services
	userService := NewUserService(userRepo)
	//Controller Handler Same thing
	userHandler := NewUserController(userService)

	//Return the module to be plugged in
	return &Module{
		userService:    userService,
		userHandler:    userHandler,
		authMiddleware: params.AuthMiddleware,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/all", m.userHandler.FindAll)

		r.Group(func(r chi.Router) {
			r.Use(m.authMiddleware)

			r.Post("/", m.userHandler.Create)
			r.Get("/{id}", m.userHandler.GetById)
			r.Delete("/{id}", m.userHandler.DeleteById)
			r.Delete("/{id}/soft", m.userHandler.SoftDeleteById)
		})

	})
}
