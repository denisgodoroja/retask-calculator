# Variables
BINARY_NAME=pack-calculator
MAIN_GO_PKG=./cmd/web/main.go
BIN_DIR=build

.PHONY: all
all: build

.PHONY: clean
clean: ## Clean up build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "Clean complete."

.PHONY: build
build: ## Build the Go binary for the local system
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_GO_PKG)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

.PHONY: run
run: ## Run the Go application locally (not in Docker)
	@echo "Starting Go backend server (local)..."
	@go run $(MAIN_GO_PKG)

.PHONY: test
test: ## Run all unit tests
	@echo "Running tests..."
	@go test ./... -v
