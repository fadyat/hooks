lint:
	golangci-lint run

up:
	docker-compose -f build/docker-compose.yml up --build api

swag:
	swag init --generalInfo cmd/hooks/main.go --output api/docs
	swag fmt

hooks:
	go run cmd/hooks/main.go

tests:
	go test -v ./test