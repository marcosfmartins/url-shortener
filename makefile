COMPOSE_FILE := docker-compose.yml

.PHNY: lint
lint:
	@echo "Running golangci-lint..."
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.2.1 golangci-lint run

.PHONY: test
test:
	mkdir -p coverage
	docker compose --profile test up test --build
	docker compose --profile test down --volumes --remove-orphans
	sed -i '/_mock.go/d' coverage/coverage.out
	go tool cover -html=coverage/coverage.out -o coverage/index.html
	rm -f coverage/coverage.out

.PHONY: run
run:
	docker compose -f $(COMPOSE_FILE) --profile local up
