package user

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler
type authenticatorMiddlewareInterface interface {
	RequiresAuthentication(next http.Handler) http.Handler
}
type anonymousMiddlewareInterface interface {
	EnsureAnonymous(next http.Handler) http.Handler
}

// Any externel dependencies (like from the core app) needed by
// this modulei
type ModuleParams struct {
	DB            *gorm.DB
	Authorisation authenticatorMiddlewareInterface
	Anonymous     anonymousMiddlewareInterface
}

type Module struct {
	userService UserServiceInterface
	userHandler *UserHandler

	// Authorise *Middleware
	Authorise authenticatorMiddlewareInterface
	Sess      anonymousMiddlewareInterface
}

// Implement the applications's Module interface
func (m *Module) Name() string {
	return "User Module"
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
		userService: userService,
		userHandler: userHandler,
		Authorise:   params.Authorisation,
		Sess:        params.Anonymous,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(browse chi.Router) {
		browse.Use(m.Sess.EnsureAnonymous) //Perhaps this should be public (no middlware)
		browse.Get("/all", m.userHandler.GetAll)
		browse.Post("/", m.userHandler.Create)

		browse.Group(func(protected chi.Router) {
			protected.Use(m.Authorise.RequiresAuthentication)

			//Define the roles required by these routes <MOVE THIS TO THE AppMiddleware
			protected.Use(m.Sess.EnsureAnonymous)
			protected.Get("/{id}", m.userHandler.GetById)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(m.Authorise.RequiresAuthentication)

			protected.Use(m.Sess.EnsureAnonymous)
			protected.Delete("/{id}", m.userHandler.DeleteById)
			protected.Delete("/{id}/soft", m.userHandler.SoftDeleteById)
		})
	})
}
