package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

type Bucket struct {
	counter uint
	buckets [256]*Bucket
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (b *Bucket) addRecursive(ip net.IP, depth int) bool {
	if depth == 0 {
		if b.counter != 0 {
			return false
		}
		b.counter++
		return true
	}
	octet := ip[0]
	remain := ip[1:]
	if b.buckets[octet] == nil {
		b.buckets[octet] = NewBucket()
	}
	ok := b.buckets[octet].addRecursive(remain, depth-1)
	if ok {
		b.counter++
	}
	return ok
}

func (b *Bucket) Add(ip net.IP) (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("IP address is nil")
	}
	ip = ip.To4()
	if ip == nil {
		return false, fmt.Errorf("Incorrect IPv4 address: %s", ip)
	}
	return b.addRecursive(ip, 4), nil
}

func (b *Bucket) Count() uint {
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
			fmt.Printf("Incorrect IP address: %s\n", line)
			continue
		}
		ok, err := root.Add(ip)
		if err != nil {
			fmt.Printf("Couldn't add %s: %v\n", line, err)
			continue
		}
		if !ok {
			fmt.Printf("Duplicate IPv4 address: %s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("Unique IP addresses count: %d\n", root.Count())
}
