package auth

import (
	"booking-platform/internal/auth/controllers"
	"booking-platform/internal/auth/repositories"
	"booking-platform/internal/auth/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo repositories.UserRepository
}

// Module Services that hold all business logic for a domain in all contexts
type Services struct {
	UserSvc *services.UserService
}

// Module Controller Handler Same thing
type Controllers struct {
	UserHTTP *controllers.UserController
}

// Module Structure it self
type Module struct {
	Repos       Repositories
	Services    Services
	Controllers Controllers
}

//Any external dependencies (like from the core app) needed by
//this module

type ModuleParams struct {
	DB *gorm.DB
}

func NewModule(params ModuleParams) (*Module, error) {

	//register repos
	repos := Repositories{
		UserRepo: *repositories.NewUserRepo(params.DB),
	}
	//register services
	services := Services{
		UserSvc: services.NewUserService(repos.UserRepo)}

	//Controller Handler Same thing
	controllers := Controllers{
		controllers.NewUserController(services.UserSvc),
	}

	//Return the module to be plugged in
	return &Module{
		Repos:       repos,
		Services:    services,
		Controllers: controllers,
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(sub chi.Router) {
		sub.Post("/", m.Controllers.UserHTTP.Create)
	})
}
