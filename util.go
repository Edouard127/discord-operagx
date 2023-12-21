package main

import (
	"context"
	"time"
)

func Ticker(ctx context.Context, every time.Duration, f func()) {
	last := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Since(last) >= every {
				f()
				last = time.Now()
			}
		}
	}
}
