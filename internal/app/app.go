// @title		Booking Platform Service
// @version		0.1
// @description	This is a backend service for a booking platform(not like Booking.com)
// @host		localhost:8080
// @BasePath	/
package app

import (
	"context"
	"fmt"

	//unused import for swag AST to include these packages
	_ "booking-platform/docs/swagger"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Module interface {
	Name() string
	RegisterRoutes(chi.Router)
}

// Making every module implement Start() and Stop() can get messy
type Startable interface {
	Start(context.Context) error
}
type Stopable interface {
	Stop(context.Context) error
}
type Application struct {
	Modules []Module
	Router  chi.Router
}

func NewApplication(modules ...Module) (*Application, error) {
	//Creating our 'root' application router
	r := chi.NewRouter()

	//Logger needs to come first before any other middleware that can modify requests
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	//Recoverer recovers from panics,logs the panic (with the backtrace),
	// then returns HTTP 500 status
	r.Use(middleware.Recoverer)
	//Swagger UI dashboard
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	app := &Application{
		Router: r,
	}
	//Modules are initialised here
	for _, module := range modules {
		err := app.RegisterModule(module)
		if err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (a *Application) RegisterModule(module Module) error {
	if module == nil {
		return fmt.Errorf("Cannot register nil modules")
	}
	a.Modules = append(a.Modules, module)
	module.RegisterRoutes(a.Router)

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	// This function is supposed to start the plugged in modules
	for _, module := range a.Modules {
		startable, ok := module.(Startable)
		if !ok {
			//module doesn't implement Start
			continue
		}
		err := startable.Start(ctx)
		if err != nil {
			return fmt.Errorf("Starting Module: %q: %w", module.Name(), err)
		}
	}
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	//stop modules like popping off of a stack
	for i := len(a.Modules) - 1; i >= 0; i-- {
		module := a.Modules[i]
		stopable, ok := module.(Stopable)
		if !ok {
			//module doesn't implement Start
			continue
		}
		err := stopable.Stop(ctx)
		if err != nil {
			return fmt.Errorf("Stopped Module %s:%w", module.Name(), err)
		}
	}
	return nil
}
