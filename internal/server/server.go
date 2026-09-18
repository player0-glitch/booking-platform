package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"booking-platform/internal/app"
	"booking-platform/internal/database"
)

type Server struct {
	port int
	db   database.Service
	app  *app.Application
}
type Params struct {
	Application *app.Application
	Database    database.Service
	// Port        int
}

func NewServer(params Params) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	// database, _ := database.New()

	//Init server
	serverInstance := &Server{
		port: port,
		// port:params.Port
		//Pass in a pointer of applications
		app: params.Application,
		db:  params.Database,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverInstance.port),
		Handler: serverInstance.Handler(),
		// Handler:      NewServer.app.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
