.PHONY: run build test clean gen tidy migrate

# Disable go.work by default for monorepo compatibility (override with GOWORK=on).
GOWORK ?= off

# Run the application
run:
	GOWORK=$(GOWORK) go run ./cmd/main.go

# Build the application
build:
	GOWORK=$(GOWORK) go build -o bin/app ./cmd/main.go

# Run tests
test:
	GOWORK=$(GOWORK) go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Generate a domain module
# Usage: make gen NAME=User FIELDS="username,email,status:enum(active|inactive)"
gen:
	soliton-gen domain $(NAME) --fields "$(FIELDS)"

# Run database migrations
migrate:
	GOWORK=$(GOWORK) go run ./cmd/migrate

# Tidy dependencies
tidy:
	GOWORK=$(GOWORK) go mod tidy
