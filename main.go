package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

const API = "https://api.discord.gx.games/v1/direct-fulfillment"
const Template = "https://discord.com/billing/partner-promotions/1180231712274387115/%s\n"

var file, _ = os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		slog.Error("could not connect to tor", slog.String("error", err.Error()))
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}

	slog.Info("connected to the tor network")

	var discord Discord
	var resp *http.Response

	ctx, fn := context.WithCancelCause(context.Background())
	defer fn(nil)

	Ticker(ctx, 5*time.Second, func() {
		resp, err = client.Post(API, "application/json", bytes.NewBuffer([]byte("{\"partnerUserId\":\"cb7f04df-8b8e-4dc8-bc20-2b0e60e211d9\"}")))
		if err != nil {
			slog.Error("could not make request to discord api", slog.String("error", err.Error()))
			fn(err)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&discord)
		if err != nil {
			slog.Error("could not decode the discord api response", slog.String("error", err.Error()))
			fn(err)
			return
		}

		slog.Info("discord api response", slog.String("token", discord.Token))
		file.Write([]byte(fmt.Sprintf(Template, discord.Token)))
	})

	file.Close()
}

func Ticker(ctx context.Context, every time.Duration, f func()) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f()
		case <-ctx.Done():
			return
		}
	}
}

type Discord struct {
	Token string `json:"token"`
}
