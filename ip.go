package ipset

import (
	"encoding/xml"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type IPSetIP struct {
	name       string
	setType    string
	setOptions []string
	created    bool
}

func NewIP(name string, setType string, create_options ...string) (*IPSetIP, error) {
	switch setType {
	case "hash:ip":
	case "bitmap:ip":
	default:
		return nil, fmt.Errorf("unsupported type [%s], supported type: hash:ip, network:ip", setType)
	}
	if len(name) > 31 {
		return nil, fmt.Errorf("max 31 character name")
	}

	s := &IPSetIP{
		name:       name,
		setType:    setType,
		setOptions: create_options,
	}
	err := s.Create()
	if err != nil {
		return nil, err
	}
	registry.Lock()
	defer registry.Unlock()
	registry.renameHooks[name] = func(n string) {
		s.name = n
	}

	return s, nil
}
func (s *IPSetIP) Create() error {
	arglist := []string{"create", "-exist", s.name, s.setType}
	arglist = append(arglist, s.setOptions...)
	cmd := exec.Command(bin, arglist...)
	return runError(cmd)
}
func (s *IPSetIP) List() ([]net.IP, error) {
	cmd := exec.Command(bin, "list", s.name, "-o", "xml")
	ipList := []net.IP{}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []net.IP{}, err
	}
	ipsXML := &Ipsets{}
	err = xml.Unmarshal(out, ipsXML)
	if err != nil {
		return []net.IP{}, err
	}
	for _, v := range ipsXML.Ipset[0].Members.Member {
		ip := net.ParseIP(v.Elem)
		if ip == nil {
			return []net.IP{}, fmt.Errorf("error parsing [%s] as IP", v.Elem)
		}
		ipList = append(ipList, ip)
	}
	return ipList, nil
}

func (s *IPSetIP) Swap(newName string) error {
	err := runError(exec.Command(bin, "swap", s.name, newName))
	if err != nil {
		return err
	}
	registry.Lock()
	defer registry.Unlock()
	oldName := s.name
	if _, ok := registry.renameHooks[oldName]; ok {
		//if both exist
		if f, ok := registry.renameHooks[newName]; ok {
			// tell new one to change name to the old one
			f(oldName)
			// swap them
			registry.Swap(oldName, newName)
			s.name = newName
		}
	} else {
		registry.renameHooks[newName] = func(n string) {
			s.name = n
		}
	}
	return nil
}

// Add an element. adding same element twice is no-op
func (s *IPSetIP) Add(ip net.IP) error {
	return runError(exec.Command(bin, "add", "-exist", s.name, ip.String()))
}

// Delete an element. Deleting nonexistend element is noop
func (s *IPSetIP) Delete(ip net.IP) error {
	return runError(exec.Command(bin, "del", "-exist", s.name, ip.String()))
}

func (s *IPSetIP) Destroy() error {
	err := runError(exec.Command(bin, "destroy", s.name))
	if err == nil {
		registry.Lock()
		defer registry.Unlock()
		delete(registry.renameHooks, s.name)
	}
	return err
}

// Exists returns whether set exists. It does not handle errors
func (s *IPSetIP) Exist() bool {
	cmd := exec.Command(bin, "list", "-name", s.name)
	out, _ := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	if outStr == s.name {
		return true
	}
	return false
}

func (s *IPSetIP) Name() string {
	return s.name
}
