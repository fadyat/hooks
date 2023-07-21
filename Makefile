lint:
	@golangci-lint run ./...

up:
	@docker-compose -f build/docker-compose.yml up --build api

up-prod:
	@docker-compose -f build/docker-compose.yml up api_registry

swag:
	swag init --generalInfo cmd/hooks/main.go --output api/docs
	swag fmt

run:
	@go run cmd/hooks/main.go

tests:
	@go test -v ./...

tagging:
	@echo "Tagging version $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

.PHONY: lint, up, swag, run, tests