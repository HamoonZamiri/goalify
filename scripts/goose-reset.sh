#!/bin/bash
# Get the directory of the script
SCRIPT_DIR=$(dirname "$0")

# Set the relative path to the migrations directory
MIGRATIONS_DIR="$SCRIPT_DIR/../backend/db/migrations/"
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=goalify password=goalify dbname=goalify host=localhost port=5432 sslmode=disable" goose -v -dir "$MIGRATIONS_DIR" reset
