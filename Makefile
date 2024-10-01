DB_URL ?= postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable

############################### Requirements ###############################
setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4

############################### Migrate ###############################
create-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$migration_name

migrate-up:
	migrate -database ${DB_URL} -path internal/database/migrations up

migrate-down:
	migrate -database ${DB_URL} -path internal/database/migrations down

############################### Sqlc ###############################
sqlc:
	docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate

############################### Docs ###############################
gen-docs:
	swag init

############################### App ###############################
run-dev:
	go run main.go

run-prod:
	go build -o url_shortener main.go
	./url_shortener