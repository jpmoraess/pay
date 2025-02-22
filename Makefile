DB_URL=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable

server:
	go run main.go

migrateup:
	migrate -path ./db/migration -database "${DB_URL}" -verbose up

migratedown:
	migrate -path ./db/migration -database "${DB_URL}" -verbose down

.PHONY: server migratup