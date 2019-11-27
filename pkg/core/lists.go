package core

import (
	"net"
	"sync"

	"github.com/yl2chen/cidranger"
)

// List is threadsafe container for Ranger type
// cidranger uses path-compressed trie for storing CIDRs and looking up for IP inclusion in stored CIDRs
// alternative is to use map with all black/white listed CIDRs and scan each value to check if given IP in every range (with O(n) complexity)
// another possible alternative is radix tree with O(k)
type List struct {
	mx     sync.RWMutex
	ranger cidranger.Ranger
}

// NewList returns instence of List
func NewList() *List {
	l := &List{
		ranger: cidranger.NewPCTrieRanger(),
	}
	return l
}

// InsertCIDR adds CIDR to list
func (l *List) InsertCIDR(cidr string) error {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	l.mx.Lock()
	defer l.mx.Unlock()
	err = l.ranger.Insert(cidranger.NewBasicRangerEntry(*net))
	if err != nil {
		return err
	}
	return nil
}

// LookUpIP will get IP and try to find in in current List
func (l *List) LookUpIP(ip string) (bool, error) {
	IP := net.ParseIP(ip)
	l.mx.RLock()
	defer l.mx.RUnlock()
	contains, err := l.ranger.Contains(IP)
	if err != nil {
		return false, err
	}
	return contains, nil
}

// DeleteCIDR will exclude CIDR from List
func (l *List) DeleteCIDR(cidr string) error {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	l.mx.Lock()
	defer l.mx.Unlock()
	_, err = l.ranger.Remove(*net)
	if err != nil {
		return err
	}
	return nil
}
