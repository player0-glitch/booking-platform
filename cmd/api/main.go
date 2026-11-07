package main

import (
	"booking-platform/internal/app"
	"booking-platform/internal/authentication"
	"booking-platform/internal/core/adapters"
	"booking-platform/internal/database"
	"booking-platform/internal/database/seeder"
	"booking-platform/internal/jwt"
	"booking-platform/internal/modules/booking"
	"booking-platform/internal/modules/user"
	"booking-platform/internal/server"
	"booking-platform/internal/sessions"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	Logger "github.com/sirupsen/logrus"
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

var (
	jwtSecret       = os.Getenv("JWT_SECRET")
	sessionHashKey  = os.Getenv("SESSION_KEY_SECRET")
	sessionBlockKey = os.Getenv("SESSION_BLOCK_SECRET")
)

func main() {
	//A global logger for debugging so we can se what file and line made the log
	Logger.SetReportCaller(true)
	//init database
	databaseContext, err := database.New()
	if err != nil {
		Logger.Fatalf("Failed To Start Database Service: %s", err)
	}
	//defer closing the connection
	defer func() {
		err := databaseContext.Close()
		if err != nil {
			Logger.Printf("Failed To Close Database Connection: %v", err)
		}
	}()
	//Initialise the modularized application in reverse order, what's on bottom
	// depeneds on what's above it
	// Building Up The Core Components
	fmt.Println("Building Up The Core Components")
	fmt.Println("Building JWT")
	jwtSvc := jwt.NewJWTService(jwtSecret)

	fmt.Println("Building Application Sessions")
	sess := sessions.NewSesssionStore(databaseContext.DB(),
		[]byte(sessionHashKey),
		[]byte(sessionBlockKey),
		database.NewGormTransactor(databaseContext.DB()),
	)
	sessStore := sessions.NewManager(sess, "booking_app_sess")
	sessMiddleware := sessions.NewSessionMiddleware(sessStore)

	fmt.Println("Registering Modules")

	authenticationModule := authentication.NewModule(authentication.ModuleParams{
		UserReaderInterface: fulfillAuthUserContract(databaseContext),
		DB:                  databaseContext.DB(),
		JWTGenerator:        jwtSvc,
		HandlerStoreManager: sessStore,
	})
	fmt.Printf("Registered Module 1 as %s\n", authenticationModule.Name())

	userModule := user.NewModule(user.ModuleParams{
		DB:            databaseContext.DB(),
		Authorisation: authenticationModule.Middleware,
		Anonymous:     sessMiddleware,
	})
	fmt.Printf("Registered  Module 2 as %s\n", userModule.Name())

	bookingModule := booking.NewModule(booking.ModuleParams{
		DB:         databaseContext.DB(),
		Transactor: database.NewGormTransactor(databaseContext.DB()),
	})
	fmt.Printf("Registered Module 3 %s\n", bookingModule.Name())

	//Run any seeder before the app starts
	runSeeders(context.Background(), databaseContext)

	app, errAppStart := app.NewApplication(
		authenticationModule,
		userModule,
		bookingModule,
		// catalogueModule,
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
		// Port:        8080,
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
// Auth module wants a UserReaderInterface with GetByEmailWithRoles and defines the contract
// UserRepository implements the interface/contract
func fulfillAuthUserContract(db database.Service) *adapters.UserAdapter {
	userRepo := user.NewUserRepo(db.DB())
	return adapters.NewUserAdapter(userRepo)
}

func runSeeders(ctx context.Context, db database.Service) {
	userRepo := user.NewUserRepo(db.DB())
	seeder := seeder.NewSeeder(db.DB(), seeder.SeederParams{
		UserRepo:    userRepo,
		Transaction: database.NewGormTransactor(db.DB()),
	})
	fmt.Print("Running Seeder")
	seeder.Seed(ctx)
}
