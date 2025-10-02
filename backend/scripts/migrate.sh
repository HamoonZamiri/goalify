go mod install && go tool goose -dir=../db/migrations postgres "$DATABASE_URL" up
