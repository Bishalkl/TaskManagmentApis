
# ————————————————————————————————————————————————
# Set app-related paths
# ————————————————————————————————————————————————
APP_NAME = auth-service
CMD_PATH = ./cmd/server

# .env support
ENV_FILE = .env

#DB connection from .env
DB_DRIVER = postgres
DB_STRING = "host=$${DB_HOST} port=$${DB_PORT} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} sslmode=disable"

#Goose migration directory
MIGRATION_DIR = ./migration


# ————————————————————————————————————————————————
# 1️⃣ Environment
# ————————————————————————————————————————————————
# load-env:
#   Load your .env file into the shell so that
#   subsequent commands pick up DB_HOST, JWT_SECRET, etc.
load-env:
	@export $$(cat $(ENV_FILE) | xargs)

# ————————————————————————————————————————————————
# 2️⃣ Run & Build
# ————————————————————————————————————————————————
# run:       Run the server with live code
# build:     Build a production binary
# fmt:       go fmt all files
run:
	@echo "🚀 Starting server"
	go run $(CMD_PATH)/main.go

build:
	@echo "📦 Building binary"
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

fmt:
	@echo "🖌️  Formatting code"
	go fmt ./...


# ————————————————————————————————————————————————
# 4️⃣ Goose Migrations
# ———————————

migrate-create:
	#make migrate-create name=create_yours_table
	@echo "✏️  Creating new migration: $(name)"
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	#make migrate-up
	@echo "⬆️  Applying migrations up"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up

migrate-down:
	#make migrate-down
	@echo "⬇️ Rolling back last migration"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down

migrate-status:
	#make migrate-status
	@echo "📋 Migration status"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) status

migrate-up-to:
	#make migrate-up-to version=2
	@echo "⬆️  Migrating up to version $(version)"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up $(version)

migrate-down-to:
	#make migrate-down-to version=2
	@echo "⬇️ Rolling down to version $(version)"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down $(version)

migrate-reset:
	#make migrate-reset
	@echo "⚠️  Resetting database (down all, up all)"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) reset

migrate-version:
	#make migrate-version
	@echo "🔢 Current DB version"
	@export $$(grep -v '^\s*#' .env | grep -v '^\s*$$' | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) version

migrate-fix:
	#make migrate-fix
	@echo "🛠️  Fixing duplicate migration numbering"
	@goose -dir $(MIGRATIONS_DIR) fix