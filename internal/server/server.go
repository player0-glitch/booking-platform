package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"booking-platform/internal/application"
	"booking-platform/internal/database"
)

type Server struct {
	port int
	db   database.Service
	app  *application.Application
}
type Params struct {
	Application *application.Application
}

func NewServer(params Params) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	database, _ := database.New()

	//Init server
	NewServer := &Server{
		port: port,
		//Pass in a pointer of applications
		app: params.Application,
		db:  database,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
		// Handler:      NewServer.app.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
