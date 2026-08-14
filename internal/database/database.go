package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	//replace sql* DB with gorm* DB
	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	// DB returns the underlying GORM database instance that should
	//only exist once I think
	DB() *gorm.DB
}

type service struct {
	// db *sql.DB
	db *gorm.DB
}

var (
	dburl = os.Getenv("DEV_DB_URL")
	// dbInstance *service
)

// Trying to prevent 2 go routines from both intialising this service
var (
	dbInstance Service
	dbOnce     sync.Once
)

func New() Service {
	// Reuse Connection
	dbOnce.Do(func() {

		db, err := gorm.Open(
			sqlite.Open(dburl), &gorm.Config{})

		// //replace sql with gorm
		// db, err := gorm.Open(sqlite.Open(dburl),
		// 	&gorm.Config{})

		if err != nil {
			// This will not be a connection error, but a DSN parse error or
			// another initialization error.
			log.Fatal(err)
		}

		//ping db to make sure it is connected once.
		//I think this is to ensure that al
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		if err := sqlDB.Ping(); err != nil {
			log.Fatal(err)
		}
		dbInstance = &service{
			db: db,
		}

	})
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	//Get the underlying *sql.DB from GORM
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: :%v", err)
		return stats
	}
	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	sqlDB, err := s.db.DB()

	if err != nil {
		return err
	}

	log.Printf("Disconnected from database: %s", dburl)
	return sqlDB.Close()
}

// DB returns the underlying GORM database instance.
func (s *service) DB() *gorm.DB {
	return s.db
}
