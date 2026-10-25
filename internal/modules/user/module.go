package user

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler

type Module struct {
	userService    UserServiceInterface
	userHandler    *UserHandler
	authMiddleware *Middleware
}

// Implement the applications's Module interface
func (m *Module) Name() string {
	return "user module"
}

// Any externel dependencies (like from the core app) needed by
// this modulei
type ModuleParams struct {
	DB             *gorm.DB
	AuthMiddleware Middleware
}

func NewModule(params ModuleParams) *Module {
	//register repos
	userRepo := NewUserRepo(params.DB)
	rolesRepo := newRolesRepo(params.DB)
	//register services. Our service takes in a hard copy
	//for sustainablility. Just point to the Implementation of choice
	userService := NewUserService(*userRepo, rolesRepo)
	//Controller Handler Same thing
	userHandler := NewUserController(userService)

	//Return the module to be plugged in
	return &Module{
		userService:    userService,
		userHandler:    userHandler,
		authMiddleware: &params.AuthMiddleware,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/all", m.userHandler.GetAll)
		r.Post("/", m.userHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(*m.authMiddleware)

			//Define the roles required by these routes
			r.Use(authorisation.RequiredRoles(
				core.RoleAdmin,
				core.RoleCustomer, /*provision guest*/
				core.RoleGuest))
			r.Get("/{id}", m.userHandler.GetById)
		})
		//Should fallback to Provisioning a Guest User
		r.Get("/{rolename}", m.userHandler.GetUsersByRole)

		r.Group(func(r chi.Router) {
			r.Use(*m.authMiddleware)

			r.Use(authorisation.RequiredRoles(core.RoleSuperAdmin, core.RoleAdmin))
			// r.Post("/", m.userHandler.Create)
			r.Delete("/{id}", m.userHandler.DeleteById)
			r.Delete("/{id}/soft", m.userHandler.SoftDeleteById)
		})
	})
}
