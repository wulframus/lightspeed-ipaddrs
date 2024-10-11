// Implements BitSet algorithm to accurately count unique IPv4 addresses from a big bulk of data.
// It is fast and uses a fixed memory size.
package ipcounter

import (
	"net"

	"github.com/bits-and-blooms/bitset"
)

const (
	storageSize = (1 << 32) // Length of storage enough for all IPv4 addresses range
)

// Unique IPv4 counter structure. Helps efficiently count unique IPv4 addresses from big bulk of data.
type UniqueIpv4Counter struct {
	storage *bitset.BitSet
	counter uint
}

// IPv4 address to UInt32 converter.
func v4ToUInt32(ip net.IP) uint32 {
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}

// NewUniqueIpv4Counter creates and returns UniqueIpv4Counter instance.
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

// Add adds an IPv4 address to storage and increases counter if it hasn't been in.
func (c *UniqueIpv4Counter) Add(ip string) {
	parsedIp := net.ParseIP(ip)
	parsedIp = parsedIp.To4()
	c.addUInt32(v4ToUInt32(parsedIp))
}

// Clear clears all stored addresses and zeros counter. It doesn't free memory.
func (c *UniqueIpv4Counter) Clear() {
	c.storage.ClearAll()
	c.counter = 0
}

// Count returns a current counter value
func (c *UniqueIpv4Counter) Count() uint {
	return c.counter
}
