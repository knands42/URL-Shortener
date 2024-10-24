include app-dev.env

DB_URL ?= postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable
DB_URL_TEST ?= postgres://postgres:postgres@localhost:5434/url_shortener_test?sslmode=disable

############################### Requirements ###############################
setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/perf/cmd/benchstat@latest
	brew install golang-migrate
	python3 -m venv ./.venv
	./.venv/bin/pip install bzt

############################### Migrate ###############################
create-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$migration_name

migrate-up:
	migrate -database ${DB_URL} -path internal/database/migrations up
	migrate -database ${DB_URL_TEST} -path internal/database/migrations up

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

stop-b-and-r:
	killall -e url_shortener	

build-and-run:
	go build -o url_shortener main.go
	./url_shortener

integration-tests:
	go test ./tests/integration -cover

performance-tests:
	go test -bench=. -benchmem ./tests/performance