PORT=80

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

pre:
	@go mod tidy
	@make swag
	@make lint
	@make tests

lint:
	@golangci-lint run ./...

up:
	@docker-compose -f build/docker-compose.yml up --build api

up-prod:
	@docker-compose -f build/docker-compose.yml up --build api_registry

swag:
	@swag init --generalInfo cmd/hooks/main.go --output api/docs
	@swag fmt

run:
	@go run cmd/hooks/main.go

tests:
	@go test -cover ./...

recreate-tag: delete-tag tagging

tagging:
	@echo "Tagging version $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

delete-tag:
	@echo "Deleting tag $(VERSION)"
	@git tag -d $(VERSION)
	@git push origin --delete $(VERSION)

ngrok:
	@ngrok http $(PORT)

.PHONY: pre, lint, up, swag, run, tests, tagging, ngrok