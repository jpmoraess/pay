DB_URL=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable

server:
	go run main.go

.PHONY: server