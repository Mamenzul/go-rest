build:
	@go build -o bin/go-rest main.go
	
run: build
	@./bin/go-rest

migration:
	@migrate create -ext sql -dir migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run migrate/main.go up

migrate-down:
	@go run migrate/main.go down
