# Core GoLang dependencies needed

- Gonertia: Inertia.js wrapper for GoLang
  go get github.com/romsar/gonertia/v3

- Chi: high performance router for http handlers (Controller)
  go get github.com/go-chi/chi/v5

- Swag: we're back to swagger baby
  go install github.com/swaggo/swag/cmd/swag@latest

- Gorilla Seessions: Manage cookies and session infrastructure
  go get github.com/gorilla/sessions
  go get github.com/gorilla/csrf

- GORM: ORM for go along with sqlite3 drivers for development
  go get gorm.io/gorm
  go get gorm.io/driver/sqlite

- Asynq: this may require redis later on though
  go get github.com/hibiken/asynq

- Redis client
  go get github.com/redis/go-redis/v9

- Air: for dev experience to have hot reload with api and worker configs
  go install github.com/air-verse/air@latest
