# psql -p 5432 -h localhost --username goalify
#!/bin/bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=goalify password=goalify dbname=goalify host=localhost port=5432 sslmode=disable" goose -v -dir "/Users/hamoonzamiri/Desktop/Laptop/Projects/goalify/backend/db/migrations/" up
