package main

import (
	"protun/server"
	"protun/vpn"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	started := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		vpn.StartSession(started)
	}()
	<-started

	server.StartProxyServer(&wg)
}
