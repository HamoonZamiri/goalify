build:
	@go build -o bin/goalify

run: build
	@./bin/goalify
unit:
	@echo "Running unit tests..."
	@go test -count=1 ./utils/... ./middleware/... ./db/...
inte:
	@echo "Running integration tests..."
	@go test -count=1 -v ./users/stores/...
	@go test -count=1 -v ./goals/stores/...
	@go test -count=1 -v ./main_test.go
test: unit inte
up:
	@docker compose up -d
down: 
	@docker compose down

