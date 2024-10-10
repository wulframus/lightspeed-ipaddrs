package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

type Tree struct {
	leftCount int
	leftIsFull bool
	leftSubTree *Tree
	rightCount int
	rightIsFull bool
	rightSubTree *Tree
}

func NewTree() *Tree {
	return &Tree{}
}

func (self *Tree) IsFull() bool {
	return self.leftIsFull && self.rightIsFull
}

func (self *Tree) Count() uint {
	return uint(self.leftCount + self.rightCount)
}

func (self *Tree) addRecursive(ip uint32, depth int) int {
	if self.IsFull() {
		return 0
	}
	ret := 0
	if depth == 0 {
		if self.leftIsFull != self.rightIsFull {
			self.leftIsFull = true
			ret = -1
		} else {
			self.rightIsFull = true
			ret = 1
		}
	} else {
		if (ip & 1) != 0 {
			if self.rightIsFull {
				return ret
			}
			if self.rightSubTree == nil {
				self.rightSubTree = NewTree()
			}
			ret = self.rightSubTree.addRecursive(ip >> 1, depth - 1)
			self.rightCount += ret
			if ret == 0 {
				self.rightIsFull = true
				self.rightSubTree = nil
			}
		} else {
			if self.leftIsFull {
				return ret
			}
			if self.leftSubTree == nil {
				self.leftSubTree = NewTree()
			}
			ret = self.leftSubTree.addRecursive(ip >> 1, depth - 1)
			self.leftCount += ret
			if ret == 0 {
				self.leftIsFull = true
				self.leftSubTree = nil
			}
		}
	}
	return ret
}

func toUint32(ip net.IP) uint32 {
	return uint32(ip[0]) << 24 | uint32(ip[1]) << 16 | uint32(ip[2]) << 8 | uint32(ip[3])
}

func (self *Tree) Add(ip net.IP) {
	self.addRecursive(toUint32(ip), 32)
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
	root := NewTree()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ip := net.ParseIP(line)
		if ip == nil {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		root.Add(ip)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("Unique IP addresses count: %d\n", root.Count())
}
