# GO-IPSET

Ipset list management. So far only implemented the ip and net hash types


```go
	ipset, err := NewIP("t_ip", "hash:ip")
	err = ipset.Add(net.ParseIP("1.2.3.5"))
	err = ipset.Delete(net.ParseIP("1.2.3.5"))
    err = ipset.Destroy()
	err = ipset.Swap("t_new_ip")
```


## Testing

Testing requires permissions to manage ipsets so in most cases that's `sudo go test -v .`. Or in a container.

