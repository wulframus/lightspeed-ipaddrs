package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

type Bucket struct {
	counter int
	buckets [256]*Bucket
}

func NewBucket() *Bucket {
	return &Bucket{}
}

/**
 * returns 1 if address added, -1 if occurs once, 0 if occurs more times
 */
func (b *Bucket) addRecursive(ip net.IP, depth int) int {
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
			b.buckets[octet] = NewBucket()
		}
		ret = b.buckets[octet].addRecursive(remain, depth-1)
		b.counter += ret
	}
	return ret
}

func (b *Bucket) Add(ip net.IP) error {
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

func (b *Bucket) Count() uint {
	return uint(b.counter)
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
			fmt.Printf("Incorrect IP address: %s\n", line)
			continue
		}
		err := root.Add(ip)
		if err != nil {
			fmt.Printf("Couldn't add %s: go run -race%v\n", line, err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("Unique IP addresses count: %d\n", root.Count())
}
