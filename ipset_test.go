package ipset

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIPSet(t *testing.T) {
	var tt IPSet
	tt, err := NewIP("t_interface", "hash:ip")
	require.NoError(t, err)
	assert.NoError(t, tt.Destroy())
	tt, err = NewNet("t_interface", "hash:net")
	require.NoError(t, err)
	assert.NoError(t, tt.Destroy())
}
