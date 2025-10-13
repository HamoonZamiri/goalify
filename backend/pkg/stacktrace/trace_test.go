package stacktrace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTrace(t *testing.T) {
	d := NewDomainStackTraceLogger("test")
	trace := d.GetTrace("TestGetTrace")
	assert.Equal(t, "domain=test: TestGetTrace", trace)
}
