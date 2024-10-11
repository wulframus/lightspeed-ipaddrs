package main

import (
	"bufio"
	"context"
	"fmt"
	"lightspeed-addrs/ipcounter"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(`Usage: %s addrs.txt

Counts unique IPv4 addresses from a big text file.

!!! The program naively believes you will provide a good file without empty lines and errors :-) !!!

There is no any piping solution. Use shell tricks, like:

$ %s <(unzip -p ip_addresses.zip)

or

$ %s <(curl -s -L https://example.com/ip_addreses.txt)
`, os.Args[0], os.Args[0], os.Args[0])
		os.Exit(0)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("%s: Couldn't open file %s: %v\n", os.Args[0], os.Args[1], err)
		os.Exit(1)
	}
	defer file.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	scanner := bufio.NewScanner(file)
	counter := ipcounter.NewUniqueIpv4Counter()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Text()
				counter.Add(line)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}
	}()

	wg.Wait()

	fmt.Printf("Unique IP addresses count: %d\n", counter.Count())
}
