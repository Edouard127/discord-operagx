package main

import (
	"context"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"time"
)

func Ticker(ctx context.Context, every time.Duration, f func() bool) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if f() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func ListenForCircuit(interval time.Duration, current *http.Client) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	Ticker(context.Background(), interval, func() bool {
		current = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			},
		}
		return false
	})
}
