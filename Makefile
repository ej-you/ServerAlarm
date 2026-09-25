go_exec="./main.go"

api-dev:
	cd ./api && go run $(go_exec)

api-lint:
	cd ./api && golangci-lint run

ntfy-dev:
	docker compose -f ./dev.docker-compose.yml up -d

ntfy-dev-down:
	docker compose -f ./dev.docker-compose.yml down

mysql:
	mysql -u test_user -p server_room
