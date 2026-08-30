package application

import (
	"booking-platform/internal/auth"
	"fmt"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Module interface {
	RegisterRoutes(chi.Router)
}
type Application struct {
	// Auth *auth.Module
	Modules []Module
	// Middleware *chi.Middlewares
	Router chi.Router
}

func NewApplication(db *gorm.DB) (*Application, error) {
	//Creating our 'root' application router
	baseRouter := chi.NewRouter()

	app := &Application{
		Router: baseRouter,
	}
	//This is how you basically plug a module into the application
	authModule, err := auth.NewModule(auth.ModuleParams{
		DB: db,
	})
	if err != nil {
		fmt.Println("Failed To Register Auth Module")
		return nil, err
	}
	//plug the module with it's routes,service and db interfaces in here
	app.RegisterModule(authModule)

	return app, nil
}

func (a *Application) RegisterModule(module Module) {
	a.Modules = append(a.Modules, module)
	module.RegisterRoutes(a.Router)
}
