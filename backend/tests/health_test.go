package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Parallel()
	res, err := http.Get(BaseURL + "/health")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
