package janitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJanitor(t *testing.T) {
	j := New(map[string]time.Duration{
		"*":   500 * time.Millisecond,
		"a/*": 100 * time.Millisecond,
	})
	d, err := j.Get("b")
	assert.NoError(t, err)
	assert.Equal(t, 500*time.Millisecond, d)
	d, err = j.Get("a/b")
	assert.NoError(t, err)
	assert.Equal(t, 100*time.Millisecond, d)
	j = New(map[string]time.Duration{
		"b/*": 200 * time.Millisecond,
	})
	d, err = j.Get("c")
	assert.Error(t, err)
	assert.Equal(t, time.Duration(0), d)
}
