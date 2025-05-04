
# ————————————————————————————————————————————————
# Set app-related paths
# ————————————————————————————————————————————————
APP_NAME = auth-service
CMD_PATH = ./cmd/server

# .env support
ENV_FILE = .env

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