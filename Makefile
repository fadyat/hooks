lint:
	golangci-lint run

up:
	docker-compose up --build api

hooks:
	cd ./cmd/hooks && go run main.go

tests:
	cd ./test && go test -v ./...