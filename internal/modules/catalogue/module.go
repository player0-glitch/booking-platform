package catalogue

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler
type MiddlewareInterface interface {
	AuthMiddlware(next http.Handler) http.Handler
	EnsureAnonymous(next http.Handler) http.Handler
}
type ModuleParams struct {
	DB                  *gorm.DB
	AuthMiddleware      Middleware
	AnonymousMiddleware Middleware
}

type Module struct {
	resourceRepo *ResourceRepository
	serivceRepo  *ServicesRepository
	resourceSvc  *ResourceService
	servicesSvc  *ServicesService
}

func (m *Module) Name() string {
	return "Catalogue Module"
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/catalogue", func(browse chi.Router) {

	})
	r.Route("/organisation", func(protected chi.Router) {

	})
}
