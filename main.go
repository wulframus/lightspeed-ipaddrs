package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

type Bucket interface {
	Add(ip net.IP) error
	IsFull() bool
	Count() int
	FullBuckets() int
	addRecursive(ip net.IP, depth int) int
}

type BucketFull struct {
	counter int
	matches int
}

type BucketFeed struct {
	counter int
	fullBucketsCounter int
	buckets [256]Bucket
}

func newBucketFull(counter int) *BucketFull {
	return &BucketFull{counter: counter}
}

func newBucketFeed() *BucketFeed{
	return &BucketFeed{}
}

func NewBucket() *BucketFeed {
	return newBucketFeed()
}


// Interface implementation for BucketFull

func (b *BucketFull) IsFull() bool {
	return true
}

func (b *BucketFull) FullBuckets() int {
	return 256
}

func (b *BucketFull) addRecursive(ip net.IP, depth int) int {
	ret := 0
	if b.matches == 0 {
		ret = 1
	} else if b.matches == 1 {
		ret = -1
	}
	b.matches++
	return ret
}

func (b *BucketFull) Add(ip net.IP) error {
	return nil // Bucket is full, nothing to add
}


func (b *BucketFull) Count() int {
	return b.counter
}

// Interface implementation for BucketFeed

func (b *BucketFeed) IsFull() bool {
	return false
}

func (b *BucketFeed) FullBuckets() int {
	return b.fullBucketsCounter
} 
/**
 * returns 1 if address added, -1 if occurs once, 0 if occurs more times
 */
func (b *BucketFeed) addRecursive(ip net.IP, depth int) int {
	ret := 0
	if depth == 0 {
		if b.counter == 0 {
			ret = 1
		} else if b.counter == 1 {
			ret = -1
		}
		b.counter++
	} else {
		octet := ip[0]
		remain := ip[1:]
		if b.buckets[octet] == nil {
			if(depth > 1) {
				b.buckets[octet] = newBucketFeed()
			} else {
				b.buckets[octet] = newBucketFull(0)
				b.fullBucketsCounter++
			}
		}
		ret = b.buckets[octet].addRecursive(remain, depth - 1)
		b.counter += ret
		if !b.buckets[octet].IsFull() && b.buckets[octet].FullBuckets() == 256 {
			counter := b.buckets[octet].Count()
			b.buckets[octet] = newBucketFull(counter)
			b.fullBucketsCounter++
		}
	}
	return ret
}

func (b *BucketFeed) Add(ip net.IP) error {
	if ip == nil {
		return fmt.Errorf("IP address is nil")
	}
	ip = ip.To4()
	if ip == nil {
		return fmt.Errorf("Incorrect IPv4 address: %s", ip)
	}
	b.addRecursive(ip, 4)
	return nil
}

func (b *BucketFeed) Count() int {
	return b.counter
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <addrs.txt>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("%s: Couldn't open file %s: %v\n", os.Args[0], os.Args[1], err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	root := NewBucket()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ip := net.ParseIP(line)
		if ip == nil {
			// fmt.Printf("Incorrect IP address: %s\n", line)
			continue
		}
		err := root.Add(ip)
		if err != nil {
			// fmt.Printf("Couldn't add %s: %v\n", line, err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("Unique IP addresses count: %d\n", root.Count())
}
