package booking

import (
	"booking-platform/internal/database"
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
	Transactor          database.TransactorInterface
	AuthMiddleware      Middleware
	AnonymousMiddleware Middleware
}

type Module struct {
	repo *Repository
	svc  *Service
}

func (m *Module) Name() string {
	return "Booking Module"
}
func NewModule(params ModuleParams) *Module {
	repo := NewRepository(params.DB)
	svc := NewService(repo, params.Transactor)
	return &Module{
		repo: repo,
		svc:  svc,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/bookings", func(browse chi.Router) {

	})
}
