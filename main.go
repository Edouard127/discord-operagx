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

var file, _ = os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
var buffer = bytes.NewBuffer([]byte("{\"partnerUserId\":\"cb7f04df-8b8e-4dc8-bc20-2b0e60e211d9\"}"))
var response Discord

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	client := NewClient()
	slog.Info("connected to the tor network")

	controller := StartController("127.0.0.1:9051")

	Ticker(context.TODO(), 750*time.Millisecond, func() {
		resp, err := client.Post(API, "application/json", buffer)
		if err != nil {
			slog.Error("could not make request to discord api", slog.String("error", err.Error()))
			controller.Signal(NewCircuit)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			slog.Error("could not decode the discord api response", slog.String("error", err.Error()))
			return
		}

		slog.Info("discord api response", slog.String("token", response.Token))
		file.Write([]byte(fmt.Sprintf(Template, response.Token)))
	})

	file.Close()
}

func StartController(addr string) *Controller {
	controller, err := NewController(addr)
	if err != nil {
		slog.Error("could not connect to the tor controller", slog.String("error", err.Error()))
		return nil
	}

	err = controller.AuthenticateNone()
	if err != nil {
		slog.Error("could not authenticate with the tor controller", slog.String("error", err.Error()))
		return nil
	}

	slog.Info("connected to the tor controller")
	return controller
}

func NewClient() *http.Client {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		slog.Error("could not connect to tor", slog.String("error", err.Error()))
		return nil
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}
}
