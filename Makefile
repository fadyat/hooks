include ./build/mk/env.mk ./build/mk/docker.mk ./build/mk/tags.mk

swag:
	@swag init --generalInfo cmd/hooks/main.go --output api/docs
	@swag fmt

fmt:
	@golangci-lint run ./... --fix

lint:
	@golangci-lint run ./...

run:
	@go run cmd/hooks/main.go

tests:
	@go test -cover ./...

pre:
	@go mod tidy
	@make swag lint tests


.PHONY: pre, lint, up, swag, run, tests, tagging, ngrok