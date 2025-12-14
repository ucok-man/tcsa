# Include .env if exists
-include .env.*

# ------------------------------------------------------------------ #
#                               HELPERS                              #
# ------------------------------------------------------------------ #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'


# ------------------------------------------------------------------ #
#                             API SCRIPT                             #
# ------------------------------------------------------------------ #
## dev: run api in development mode
.PHONY: dev
dev:
	@air -c .air.toml

## build: build the cmd/api application
.PHONY: build
build:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api/

## start: start build artifact in bin/
.PHONY: start
start:
	@echo 'starting bin/api...'
	@./bin/api

## test: run all test verbose with coverage
.PHONY: test
test:
	@echo 'runing test...'
	@go test -v -cover ./...

## test/doc: run all test with gotestdox
.PHONY: test/doc
test/doc:
	@echo 'runing test...'
	@gotestdox -v -cover ./...

## swag: generate swagger definition
.PHONY: swag
swag:
	@echo 'generating swagger docs...'
	@swag init -g cmd/api/main.go  -o ./cmd/api/docs

# ------------------------------------------------------------------ #
#                          Compose Script                            #
# ------------------------------------------------------------------ #

## compose/up: run api and database using docker compose
.PHONY: compose/up
compose/up:
	@read -p "Do you want to run migration script first? [y/N] " ans; \
	if echo "$$ans" | grep -iq '^y$$'; then \
		$(MAKE) db/up; \
		$(MAKE) db/wait; \
		$(MAKE) migrate/up; \
	fi; \
	docker compose up
	
## compose/clear: clear api and database container artifact
.PHONY: compose/clear
compose/clear:
	@echo 'Cleaning up artifact...'
	@docker compose down -v

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## tidy: tidy and vendor module dependencies, and format all .go files
.PHONY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor
	@echo 'Formatting .go files...'
	go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies...'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	CGO_ENABLED=1 go test -race -vet=off -cover ./...

# ------------------------------------------------------------------ #
#                          Database Script                           #
# ------------------------------------------------------------------ #
## db/up: run database container instance
.PHONY: db/up
db/up:
	@echo 'Starting database...'
	@docker compose up -d tcsa_postgres

## db/down: remove database container instance
.PHONY: db/down
db/down:
	@echo 'Removing database...'
	@docker compose down tcsa_postgres

## db/clear: clear database container instance include volume
.PHONY: db/clear
db/clear:
	@echo 'Cleaning up database...'
	@docker compose down -v tcsa_postgres

## db/wait: wait until database is ready
.PHONY: db/wait
db/wait:
	@echo "Waiting for database to be ready..."
	@until docker compose exec -T tcsa_postgres pg_isready \
		-U ${POSTGRES_USER} \
		-d ${POSTGRES_DB} >/dev/null 2>&1; do \
		sleep 2; \
	done
	@echo "Database is ready."

# ------------------------------------------------------------------ #
#                          Migration Script                          #
# ------------------------------------------------------------------ #

## migrate/new: create new migration
.PHONY: migrate/new
migrate/new:
	@echo -n "Enter migration name: "; \
	read migration_name; \
	migration_name=$$(echo "$$migration_name" | tr ' ' '_'); \
	if [ -z "$$migration_name" ]; then \
		echo "\nError: Migration name cannot be empty." >&2; \
		exit 1; \
	fi; \
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${TCSA_DB_DSN} goose --dir ./migrations create $$migration_name sql

## migrate/up: apply all migration to latest
.PHONY: migrate/up
migrate/up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${TCSA_DB_DSN} goose --dir ./migrations up

## migrate/reset: roll back all migration
.PHONY: migrate/down
migrate/down:
	@read -p "Are you sure you want to reset the DB? [y/N] " ans; \
	if echo "$$ans" | grep -iq '^y$$'; then \
		GOOSE_DRIVER=postgres GOOSE_DBSTRING=${TCSA_DB_DSN} goose --dir ./migrations reset; \
	fi

## migrate/version: show current version applied migration
.PHONY: migrate/version
migrate/version:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${TCSA_DB_DSN} goose --dir ./migrations version

## migrate/status: dump the migration status for the current DB
.PHONY: migrate/status
migrate/status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${TCSA_DB_DSN} goose --dir ./migrations status

