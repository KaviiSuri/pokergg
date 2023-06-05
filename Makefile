build:
	@go build -o ./tmp/pokergg

# Path: Makefile
# Run the application
run:
	@go run main.go

# Path: Makefile
# Run the application with hot reload
dev:
	@air

setup:
	@go install github.com/cosmtrek/air@latest
