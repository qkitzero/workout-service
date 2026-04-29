-include .env

.PHONY: test lint lint-fix \
	proto-format proto-lint proto-gen \
	mock-gen \
	migrate-up migrate-up-one migrate-down migrate-reset migrate-create migrate-status

test:
	mkdir -p tmp
	go test -cover ./internal/... -coverprofile=./tmp/cover.out
	go tool cover -func=./tmp/cover.out | tail -n 1
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	open ./tmp/cover.html

lint:
	go tool golangci-lint run ./...

lint-fix:
	go tool golangci-lint run --fix ./...

proto-format:
	buf format -w

proto-lint:
	buf lint

proto-gen:
	buf generate

mock-gen:
	go tool mockgen -source=internal/domain/set/set.go -destination=mocks/domain/set/mock_set.go -package=mocks
	go tool mockgen -source=internal/domain/set/repository.go -destination=mocks/domain/set/mock_repository.go -package=mocks
	go tool mockgen -source=internal/domain/workout/workout.go -destination=mocks/domain/workout/mock_workout.go -package=mocks
	go tool mockgen -source=internal/domain/workout/repository.go -destination=mocks/domain/workout/mock_repository.go -package=mocks
	go tool mockgen -source=internal/domain/exercise/exercise.go -destination=mocks/domain/exercise/mock_exercise.go -package=mocks
	go tool mockgen -source=internal/domain/exercise/repository.go -destination=mocks/domain/exercise/mock_repository.go -package=mocks
	go tool mockgen -source=internal/domain/muscle/muscle.go -destination=mocks/domain/muscle/mock_muscle.go -package=mocks
	go tool mockgen -source=internal/domain/muscle/repository.go -destination=mocks/domain/muscle/mock_repository.go -package=mocks
	go tool mockgen -source=internal/application/set/usecase.go -destination=mocks/application/set/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/workout/usecase.go -destination=mocks/application/workout/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/exercise/usecase.go -destination=mocks/application/exercise/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/muscle/usecase.go -destination=mocks/application/muscle/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/auth/service.go -destination=mocks/application/auth/mock_service.go -package=mocks
	go tool mockgen -destination=mocks/external/auth/v1/mock_client.go -package=mocks github.com/qkitzero/auth-service/gen/go/auth/v1 AuthServiceClient

MIGRATIONS_DIR=internal/infrastructure/db/migrations
MIGRATE=migrate -source file://$(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_HOST_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

migrate-up:
	$(MIGRATE) up

migrate-up-one:
	$(MIGRATE) up 1

migrate-down:
	$(MIGRATE) down 1

migrate-reset:
	$(MIGRATE) drop -f

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

migrate-status:
	$(MIGRATE) version
