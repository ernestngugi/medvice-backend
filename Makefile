include .env

# make migration name=CreateUsers
migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down

up:
	docker-compose up --remove-orphans

ps:
	docker-compose ps

down:
	docker-compose down --remove-orphans

test:
	DATABASE_URL=postgres://medvice:password@localhost:5433/medvice?sslmode=disable go test ./...
