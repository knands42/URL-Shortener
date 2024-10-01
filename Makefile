DB_URL=postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable

create-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$migration_name

run-migrations:
	migrate -database ${DB_URL} -path internal/database/migrations up
