.PHONY: run build test clean gen

# Run the application
run:
	go run ./cmd/main.go

# Build the application
build:
	go build -o bin/app ./cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Generate a domain module
# Usage: make gen NAME=User FIELDS="username,email,status:enum(active|inactive)"
gen:
	soliton-gen domain $(NAME) --fields "$(FIELDS)"

# Tidy dependencies
tidy:
	go mod tidy
