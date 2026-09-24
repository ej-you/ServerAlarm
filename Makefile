go_exec="./api/main.go"

api-dev:
	go run $(go_exec)

api-lint:
	cd ./api && golangci-lint run

dev:
	docker compose -f ./docker-compose.yml up --build
