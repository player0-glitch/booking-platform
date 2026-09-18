package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"booking-platform/internal/app"
	"booking-platform/internal/database"
	"booking-platform/internal/server"
	"booking-platform/internal/user"
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
	//init database
	databaseContext, err := database.New()
	if err != nil {
		log.Fatalf("Failed To Start Database Service: %s", err)
	}
	//defer closing the connection
	defer func() {
		err := databaseContext.Close()
		if err != nil {
			fmt.Printf("Failed To Close Database Connection: %v", err)
		}
	}()

	//Initialise the modularized application
	fmt.Println("Registered Module user")
	userModule, errInitModule := user.NewModule(user.ModuleParams{
		DB: databaseContext.DB(),
	})
	if errInitModule != nil {
		log.Fatalf("%s", errInitModule.Error())
	}
	app, errAppStart := app.NewApplication(
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
		// Port:        8080,
	})
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, app, done)

	fmt.Println("Starting Server ....")
	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	fmt.Println("Wait for the graceful shutdown to complete.")
	<-done
	log.Println("Graceful shutdown complete.")
}
