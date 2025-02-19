## SQLC ##

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc init

sqlc generate

## Migrate ##

migrate create -ext sql -dir <directory> -seq <name>
Ex: migrate create -ext sql -dir db/migration -seq init_schema

up: migrate -path ./db/migration -database "database_url" -verbose up
down: migrate -path ./db/migration -database "database_url" -verbose down