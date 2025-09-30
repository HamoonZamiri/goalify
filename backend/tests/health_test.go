package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Parallel()
	res, err := http.Get(BASE_URL + "/health")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}