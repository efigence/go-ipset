package ipset

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestNewNet(t *testing.T) {
	ipset, err := NewNet("t_net", "hash:net")
	require.NoError(t, err)
	_, ip1, _ := net.ParseCIDR("1.2.3.4/24")
	_, ip2, _ := net.ParseCIDR("1.2.4.4/24")
	err = ipset.Add(ip2)
	assert.NoError(t, err)
	list, err := ipset.List()
	require.Len(t, list, 1)
	assert.Equal(t, ip2.String(), list[0].String())
	err = ipset.Add(ip1)
	assert.NoError(t, err)
	// adding same IP is noop
	err = ipset.Add(ip1)
	err = ipset.Delete(ip2)
	// deleting same IP is noop
	err = ipset.Delete(ip2)
	assert.NoError(t, err)
	list, err = ipset.List()
	assert.NoError(t, err)
	assert.Equal(t, []*net.IPNet{
		ip1,
	}, list)
	assert.True(t, ipset.Exist())
	ipset.Destroy()
}
func TestNewNetSwap(t *testing.T) {
	_, ip1, _ := net.ParseCIDR("1.2.3.4/24")
	ipset, err := NewNet("t_net", "hash:net")
	require.NoError(t, err)
	assert.NoError(t, ipset.Add(ip1))
	assert.Error(t, ipset.Swap("t_nonexistent"))
	ipset2, err := NewNet("t_net_new", "hash:net")
	assert.NoError(t, err)
	err = ipset.Swap("t_net_new")
	assert.NoError(t, err)
	assert.Equal(t, "t_net_new", ipset.Name())
	assert.Equal(t, "t_net", ipset2.Name())
	assert.True(t, ipset.Exist())
	assert.NoError(t, ipset.Destroy())
	assert.False(t, ipset.Exist())
	assert.NoError(t, ipset2.Destroy())
}
