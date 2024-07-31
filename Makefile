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
	godo test -- -e .env.local.test

test-lite:
	godo test-lite -- -e .env.local.test

coverage:
	godo coverage -- -e .env.local.test

server:
	gow run cmd/medv-api/main.go -e .env
