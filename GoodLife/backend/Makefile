DATABASE_URL ?= postgres://goodlife:goodlife@localhost:5432/goodlife?sslmode=disable

.PHONY: run migrate-up migrate-down

run:
	go run ./cmd/server

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down
