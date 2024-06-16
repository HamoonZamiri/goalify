build:
	@go build -o bin/goalify

run: build
	@./bin/goalify

test:
	@go test -count=1 ./... 

up:
	@docker-compose up -d
down: 
	@docker-compose down

