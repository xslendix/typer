package main

// https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}
	return i
}

func ScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		} else {
			// fmt.Println(port, "closed")
		}
		return
	}

	conn.Close()
	fmt.Println(ip, port, "open")
}

func StartScan(timeout time.Duration, lock *semaphore.Weighted) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for ip_addr := 0; ip_addr <= 254; ip_addr++ {
		ip := "192.168.0." + string(ip_addr)
		wg.Add(1)
		lock.Acquire(context.TODO(), 1)
		go func(port int) {
			defer lock.Release(1)
			defer wg.Done()
			ScanPort(ip, port_default, timeout)
		}(port_default)
	}
}
