package main

import (
	"booking-platform/internal/app"
	"booking-platform/internal/auth"
	"booking-platform/internal/core/adapters"
	"booking-platform/internal/database"
	"booking-platform/internal/modules/user"
	"booking-platform/internal/server"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/sirupsen/logrus"
)

func gracefulShutdown(apiServer *http.Server, application *app.Application, done chan<- bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	if err := application.Stop(ctx); err != nil {
		log.Printf("Failed to stop application modules: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	//A global logger for debugging so we can se what file and line made the log
	logger.SetReportCaller(true)
	//init database
	databaseContext, err := database.New()
	if err != nil {
		log.Fatalf("Failed To Start Database Service: %s", err)
	}
	//defer closing the connection
	defer func() {
		err := databaseContext.Close()
		if err != nil {
			logger.Printf("Failed To Close Database Connection: %v", err)
		}
	}()
	jwtSecret := os.Getenv("SESSION_STORE_KEY")
	//Initialise the modularized application
	fmt.Println("Registering Modules")
	authModule := auth.NewModule(auth.ModuleParams{
		JwtSecret:   jwtSecret,
		AuthService: fulfillAuthUserContract(jwtSecret, databaseContext),
	})
	userModule := user.NewModule(user.ModuleParams{
		DB:             databaseContext.DB(),
		AuthMiddleware: authModule.AuthMiddlware,
	})

	app, errAppStart := app.NewApplication(
		authModule,
		userModule,
	)
	if errAppStart != nil {
		log.Fatalf("%s", errAppStart.Error())
	}
	//Starting the application with context
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed To Start Application: %s", err.Error())
	}

	server := server.NewServer(server.Params{
		Application: app,
		Database:    databaseContext,
		Port:        8080,
	})
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, app, done)

	fmt.Println("Starting Server ....")
	fmt.Println(`
   ______     ______        ______     ______   __    
  /\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
  \ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
   \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  	\/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)
	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	<-done
	log.Println("Graceful shutdown complete.")
}

// Find a better way to bind these contracts
// Fullfil the contract that the Auth & User module agree on.
// Auth module wants a UserReader with GetByEmailWithRoles and defines the contract
// UserRepository implements the interface/contract
func fulfillAuthUserContract(jwtSecret string, db database.Service) *auth.AuthService {
	userRepo := user.NewUserRepo(db.DB())
	userReader := adapters.NewUserAdapter(userRepo)
	return auth.NewAuthService(userReader, []byte(jwtSecret))
}
