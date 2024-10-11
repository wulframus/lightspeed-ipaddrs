package ipcounter

import (
	"net"
	"github.com/bits-and-blooms/bitset"
)

const (
	storageSize = (1 << 32) // Length of storage enough for all IPv4 addresses range 
)

type UniqueIpv4Counter struct {
	storage *bitset.BitSet
	counter uint
}

func v4ToUInt32(ip net.IP) uint32 {
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}

func NewUniqueIpv4Counter() *UniqueIpv4Counter {
	return &UniqueIpv4Counter{
		storage: bitset.New(storageSize),
		counter: 0,
	}
}

func (c *UniqueIpv4Counter) addUInt32(i uint32) {
	if !c.storage.Test(uint(i)) {
		c.storage.Set(uint(i))
		c.counter++
	}
}

func (c *UniqueIpv4Counter) Add(ip string) {
	parsedIp := net.ParseIP(ip)
	parsedIp = parsedIp.To4()
	c.addUInt32(v4ToUInt32(parsedIp))
}

func (c *UniqueIpv4Counter) Clear() {
	c.storage.ClearAll()
	c.counter = 0
}

func (c *UniqueIpv4Counter) Count() uint {
	return c.counter
}

