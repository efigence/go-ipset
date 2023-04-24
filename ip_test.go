package ipset

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestNewIP(t *testing.T) {
	ipset, err := NewIP("t_ip", "hash:ip")
	require.NoError(t, err)
	err = ipset.Add(net.ParseIP("1.2.3.5"))
	assert.NoError(t, err)
	list, err := ipset.List()
	assert.Equal(t, []net.IP{
		net.ParseIP("1.2.3.5"),
	}, list)
	err = ipset.Add(net.ParseIP("1.2.3.4"))
	assert.NoError(t, err)
	// adding same IP is noop
	err = ipset.Add(net.ParseIP("1.2.3.4"))
	err = ipset.Delete(net.ParseIP("1.2.3.5"))
	// deleting same IP is noop
	err = ipset.Delete(net.ParseIP("1.2.3.5"))
	assert.NoError(t, err)
	list, err = ipset.List()
	assert.NoError(t, err)
	assert.Equal(t, []net.IP{
		net.ParseIP("1.2.3.4"),
	}, list)
	assert.True(t, ipset.Exist())
	ipset.Destroy()
}
func TestNewIPSwap(t *testing.T) {
	ipset, err := NewIP("t_ip", "hash:ip")
	require.NoError(t, err)
	assert.NoError(t, ipset.Add(net.ParseIP("1.2.3.5")))
	assert.Error(t, ipset.Swap("t_nonexistent"))
	ipset2, err := NewIP("t_new_ip", "hash:ip")
	assert.NoError(t, err)
	err = ipset.Swap("t_new_ip")
	assert.NoError(t, err)
	assert.Equal(t, "t_new_ip", ipset.Name())
	assert.Equal(t, "t_ip", ipset2.Name())
	assert.True(t, ipset.Exist())
	assert.NoError(t, ipset.Destroy())
	assert.False(t, ipset.Exist())
	assert.NoError(t, ipset2.Destroy())
}
