#Include environment variables
-include .env
.PHONY: all build setup dev api worker frontend worker redis test \
fmt lint redis-install debug

debug:
	@echo "MAKE IS RUNNING IN $(CURDIR)"
# Build the application
all: build test 

# consider threading this
build:
	@echo "Building appllication..."
	@make api
	@make worker
	@make frontend
	@make redis

setup:
	go mod download
	cd frontend && npm i 
	@echo "Building frontend..."
	@make redis-install

dev: 
	@make redis-install
	@echo "Running dev mode in 3 different thread"
	@make -j 3 api worker frontend

# this air commands hot reloads the api
api:
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
		echo 'Watching APIs of the project...'; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air -c .air.toml; \
                echo "Watching API...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

worker:
	@if command -v air > /dev/null; then \
		air -c .air.worker.toml; \
		echo 'Watching workers of the project...'; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air -c .air.worker.toml; \
                echo "Watching Workers..."; \
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
       fi

frontend:
	cd frontend && npm run dev

redis:
	@if redis-cli ping >/dev/null 2>&1; then \
		echo "Redis is already running..."; \
	else \
		echo 'Starting redis and by-passing homebrew'; \
		redis-server; \
	fi

# Sqlite3 is used in the development environment. 
redis-install:
	@if command -v redis-server > /dev/null; then \
		echo 'Redis Installed. Not needed for dev environment'; \
	else \
		echo 'Install redis for prod or when doing concurrency'; \
	fi

# Format everything go can find
fmt:
	go fmt ./...

# performs linting
lint:
	pwd
	go vet ./... 

clean:
	@echo "Cleaning..."
	@rm -r tmp
	@rm -f main

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# All migration files live in one migrations/ directory at root level 
# each module still has ownership of it's migrations as per the naming convention
# 000001-user-create-user.up.sql
# Notice the following;
# 000001 -> migration order
# user -> module/domain
# create-> migration action i guess i don't know
.PHONY: migrate-status migrate-up migrate-down 	migration migrate-force
DB_URL = sqlite3://storage/go_db.sqlite
MIGRATIONS_DIR ?= ./migrations/

migration-create:
	echo "$(MIGRATIONS_DIR)"
	@read -p "migration name: " name; \
	goose  create "$$name" sql

migration-status:
	@goose status
migration-up:
	@goose up
migration-down:
	@if [-z "$(version)"]; then echo "Error: Provide a 'version'. Run goose status to find versions"; exit 1; fi
	@goose down-to "$(version)"
migration-down-by-one:
	@goose down
migration-fix:
	@goose fix
migration-reset:
	@goose reset
