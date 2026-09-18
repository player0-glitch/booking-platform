package user

import (
	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Module Structure it self
type Module struct {
	userService *UserService
	// Controllers Controllers
	userHandler *UserController
}

// Implement the applications's Module interface
func (m *Module) Name() string {
	return "user"
}

// Any external dependencies (like from the core app) needed by
// this module
type ModuleParams struct {
	DB *gorm.DB
}

func NewModule(params ModuleParams) (*Module, error) {
	//register repos
	userRepo := NewUserRepo(params.DB)
	//register services
	userService := NewUserService(userRepo)
	//Controller Handler Same thing
	userHandler := NewUserController(userService)

	//Return the module to be plugged in
	return &Module{
		userService: userService,
		userHandler: userHandler,
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", m.userHandler.Create)
		r.Get("/{id}", m.userHandler.GetById)
		r.Delete("/{id}", m.userHandler.DeleteById)
		r.Delete("/{id}/soft", m.userHandler.SoftDeleteById)
		r.Get("/all", m.userHandler.FindAll)
	})
}
