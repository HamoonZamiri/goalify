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
	assert.Equal(t, "UPDATE goals SET name = $1, description = $2, completed = $3 WHERE id = $4", query)
}

func TestBuildSelectQuery(t *testing.T) {
	table := "goals"
	filters := map[string]interface{}{
		"id":         uuid.New(),
		"completed":  true,
		"created_at": "2021-01-01",
	}
	query, _ := BuildSelectQuery(table, filters)
	assert.Equal(t, "SELECT * FROM goals WHERE completed = $1 AND created_at = $2 AND id = $3", query)
}
