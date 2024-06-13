package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpdateQuery(t *testing.T) {
	table := "goals"
	updates := map[string]interface{}{
		"name":        "New Name",
		"description": "New Description",
		"completed":   true,
	}
	id := uuid.New()
	query, _ := BuildUpdateQuery(table, updates, id)
	assert.Equal(t, query,
		"UPDATE goals SET name = $1, description = $2, completed = $3 WHERE id = $4")
}

func TestBuildSelectQuery(t *testing.T) {
	table := "goals"
	filters := map[string]interface{}{
		"id":         uuid.New(),
		"completed":  true,
		"created_at": "2021-01-01",
	}
	query, _ := BuildSelectQuery(table, filters)
	assert.Equal(t, query,
		"SELECT * FROM goals WHERE id = $1 AND completed = $2 AND created_at = $3")
}
