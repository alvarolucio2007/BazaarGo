.PHONY: migrateup migratedown migration test
MIGRATIONS_PATH= ./database/migrations/
DB_URL=postgres://admin:root@localhost/scholarly?sslmode=disable
migration:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
migrateup:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose down
test:
	go test -v -cover -short ./...
run:
	go run ./...
