lint:
	golangci-lint run

up:
	docker-compose -f build/docker-compose.yml up --build api

swag:
	swag init --generalInfo cmd/hooks/main.go --output api/docs
	swag fmt

run:
	go run cmd/hooks/main.go

tests:
	go test -v ./...
