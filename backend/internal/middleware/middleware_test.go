package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIdFromHeader(t *testing.T) {
	req, err := http.NewRequest("PUT", "localhost:8000/api", nil)
	require.Nil(t, err)
	req.Header.Set("user_id", "123")

	id, err := GetIDFromHeader(req)
	require.Nil(t, err)
	assert.Equal(t, "123", id)
}
