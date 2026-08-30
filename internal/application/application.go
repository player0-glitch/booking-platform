package application

import (
	"booking-platform/internal/auth"
	"fmt"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Application struct {
	Auth *auth.Module
	// Middleware *chi.Middlewares
	Router chi.Router
}

func NewApplication(db *gorm.DB) (*Application, error) {
	//Creating our 'root' application router
	baseRouter := chi.NewRouter()

	//This is how you basically plug a module into the application
	authModule, authErr := auth.NewModule(auth.ModuleParams{
		DB: db,
	})
	if authErr != nil {
		fmt.Println("Failed To Register Auth Module")
		return nil, authErr
	}
	//plug the modules routes in here
	authModule.RegisterRoutes(baseRouter)

	return &Application{
		Auth: authModule,
		// Middleware: &baseMuxMiddleware, // we're using chi as our router and middleware
		Router: baseRouter,
	}, nil
}
