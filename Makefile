include .env

MIGRATIONS_FOLDER=./packages/db/migrations

migrate-down:
	migrate -database ${DB_URL} -path $(MIGRATIONS_FOLDER) down

migrate-up:
	migrate -database ${DB_URL} -path $(MIGRATIONS_FOLDER) up

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) -seq $(NAME) 

db-psql-connect:
	psql ${DB_URL}

swagger-gen:
	cd apps/server && swag init -g server.go

server-run:
	go run ./apps/server

server-dev:
	cd ./apps/server/ && air
