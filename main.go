package main

import (
	"rando_proxy/proxy_server"
	"rando_proxy/vpn"
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

	proxy_server.StartProxyServer(&wg)
}
