package shepherd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShepherd(t *testing.T) {
	j := NewLater(map[string]time.Duration{
		"*":   500 * time.Millisecond,
		"a/*": 100 * time.Millisecond,
	})
	d, err := j.Get("b")
	assert.NoError(t, err)
	assert.Equal(t, 500*time.Millisecond, d)
	d, err = j.Get("a/b")
	assert.NoError(t, err)
	assert.Equal(t, 100*time.Millisecond, d)
}
