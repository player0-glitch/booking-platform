package user

import (
	"booking-platform/internal/core"
	"booking-platform/internal/core/authorisation"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler
type MiddlewareInterface interface {
	AuthMiddlware(next http.Handler) http.Handler
	EnsureAnonymous(next http.Handler) http.Handler
}
type Module struct {
	userService    UserServiceInterface
	userHandler    *UserHandler
	authMiddleware *Middleware
	anonMiddleware *Middleware
}

// Implement the applications's Module interface
func (m *Module) Name() string {
	return "User Module"
}

// Any externel dependencies (like from the core app) needed by
// this modulei
type ModuleParams struct {
	DB                  *gorm.DB
	AuthMiddleware      Middleware
	AnonymousMiddleware Middleware
}

func NewModule(params ModuleParams) *Module {
	//register repos
	userRepo := NewUserRepo(params.DB)
	//register services. Our service takes in a hard copy
	//for sustainablility. Just point to the Implementation of choice
	userService := NewUserService(*userRepo)
	//Controller Handler Same thing
	userHandler := NewUserController(userService)

	//Return the module to be plugged in
	return &Module{
		userService:    userService,
		userHandler:    userHandler,
		authMiddleware: &params.AuthMiddleware,
		anonMiddleware: &params.AnonymousMiddleware,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(browse chi.Router) {
		browse.Use(*m.anonMiddleware) //Perhaps this should be public (no middlware)
		browse.Get("/all", m.userHandler.GetAll)
		browse.Post("/", m.userHandler.Create)

		browse.Group(func(protected chi.Router) {
			protected.Use(*m.authMiddleware)

			//Define the roles required by these routes <MOVE THIS TO THE AppMiddleware
			protected.Use(authorisation.RequiredRoles(
				core.RoleAdmin,
				core.RoleCustomer, /*provision guest*/
			))
			protected.Get("/{id}", m.userHandler.GetById)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(*m.authMiddleware)

			protected.Use(authorisation.RequiredRoles(core.RoleSuperAdmin, core.RoleAdmin))
			protected.Delete("/{id}", m.userHandler.DeleteById)
			protected.Delete("/{id}/soft", m.userHandler.SoftDeleteById)
		})
	})
}
