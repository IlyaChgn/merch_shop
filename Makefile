generate:
	@echo "Running go:generate"
	go generate ./...

cover:
	@echo "Running tests with coverage"
	docker compose -f docker-compose.test.yml up -d
	@until docker inspect --format='{{.State.Health.Status}}' merch_shop-postgres-1 | grep -q 'healthy'; do sleep 1; done
	go test -covermode=count -coverpkg=./... ./... -coverprofile=coverage.out
	docker compose -f docker-compose.test.yml down
	go tool cover -func=coverage.out

test:
	@echo "Running tests"
	docker compose -f docker-compose.test.yml up -d
	@until docker inspect --format='{{.State.Health.Status}}' merch_shop-postgres-1 | grep -q 'healthy'; do sleep 1; done
	go test ./... -v
	docker compose -f docker-compose.test.yml down

run:
	@echo "Running project"
	go run ./cmd/app/main.go
