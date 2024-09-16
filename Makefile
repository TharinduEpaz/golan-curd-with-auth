build:
	@go build -o bin/auth-crud cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/auth-crud
