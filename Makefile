.PHONY: migrateup migratedown migration test generate
MIGRATIONS_PATH= ./database/migrations/
DB_URL = postgresql://admin:root@localhost:5432/bazaargo?sslmode=disable
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
generate:
	oapi-codegen -config oapi-codegen.yaml api/openapi.yaml
