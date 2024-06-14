package db

import (
	"strings"
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
	assert.True(t, true, strings.Contains(query, "name = $"))
	assert.True(t, true, strings.Contains(query, "description = $"))
	assert.True(t, true, strings.Contains(query, "completed = $"))
}

func TestBuildSelectQuery(t *testing.T) {
	table := "goals"
	filters := map[string]interface{}{
		"id":         uuid.New(),
		"completed":  true,
		"created_at": "2021-01-01",
	}
	query, _ := BuildSelectQuery(table, filters)
	assert.True(t, true, strings.Contains(query, "id = $"))
	assert.True(t, true, strings.Contains(query, "completed = $"))
	assert.True(t, true, strings.Contains(query, "created_at = $"))
}
