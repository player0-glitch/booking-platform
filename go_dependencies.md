# Core GoLang dependencies needed

- Gonertia: Inertia.js wrapper for GoLang
  go get github.com/romsar/gonertia/v3
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

-Air: for dev experience to have hot reload
go install github.com/air-verse/air@latest
