## How to run the application?
```
- Goose install go get -u github.com/pressly/goose/cmd/goose
- gow install go get -u github.com/mitranim/gow
- godo install go get -u gopkg.in/godo.v2/cmd/godo
- create local env file cp .env.sample .env
- create local test env file cp .env.local.test.sample .env.local.test
- update .env file accordingly
- build dependecncies:-
 1. make up

```

Once the containers are up and ready:- 

1. run database migrations once and subsequently after adding additional migration files. 
```
make migrate
```

2. run the application by:
```
make server
```

## How to run unit tests:
Unit tests are dependent on the docker containers used for development and which should be already running.
To run the test either by:-
 1. `make test` (This spins up a new test database docker container and migrations will be run) 
 2. `make test-lite` (This does not run migrations and uses the existing database tables.)

