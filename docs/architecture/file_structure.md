# File Structure

The internal structure of the backend holding the modules of the monolith

```
internal/
├── auth/
│ ├── handlers/
│ │ └── user_handler.go
│ ├── models/
│ │ ├── user.go
│ │ ├── role.go
│ │ └── permission.go
│ ├── repositories/
│ │ └── user_repository.go
│ ├── services/
│ │ └── user_service.go
│ ├── migrate.go        -> handles ordered migrations and schema registrations
│ └── module.go         -> dependency wiring for a module. Makes it pluggable
├── database/
│   │   ├── database.go -> server config with GORM wrapper
│   │   └── migrator.go -> migrator for graceful migration handling and versioning
```
