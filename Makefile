lint:
	golangci-lint run

up:
	docker-compose -f build/docker-compose.yml up --build api

hooks:
	go run cmd/hooks/main.go

tests:
	go test -v  ./...