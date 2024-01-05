package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	// Ip - host machine ip address
	Ip string

	// lock - number of goroutines we may run simultaneously
	lock *semaphore.Weighted
}

func (p *PortScanner) scanPort(port int, d time.Duration) {
	target := fmt.Sprintf("%s:%d", p.Ip, port)
	conn, err := net.DialTimeout("tcp", target, d)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(d)
			p.scanPort(port, d)
		} else {
			fmt.Println(port, "closed")
		}
		return
	}

	conn.Close()
	fmt.Println(port, "open")
}

func (p *PortScanner) Start(f, l int, d time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		p.lock.Acquire(context.TODO(), 1)

		wg.Add(1)
		go func(port int) {
			defer p.lock.Release(1)
			defer wg.Done()
			p.scanPort(port, d)
		}(port)
	}
}

func main() {
	ps := &PortScanner{
		Ip:   "127.0.0.1",
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(1, 65535, 500*time.Millisecond)
}
