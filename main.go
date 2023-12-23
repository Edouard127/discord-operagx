package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Edouard127/controller"
	"golang.org/x/net/proxy"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

var file, _ = os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	client := NewClient()
	slog.Info("connected to the tor network")

	c := StartController("127.0.0.1:9051")

	Ticker(context.TODO(), 5*time.Millisecond, func() {
		resp, err := client.Post(API, "application/json", bytes.NewBuffer([]byte("{\"partnerUserId\":\""+UserID+"\"}")))
		if err != nil {
			slog.Error("could not make request to discord api", slog.String("error", err.Error()))
			return
		}

		if resp.StatusCode != http.StatusOK {
			slog.Error("discord api returned a non-200 status code", slog.String("status", resp.Status))
			err = c.Signal(controller.NewCircuit)
			if err != nil {
				slog.Error("could not signal the tor controller", slog.String("error", err.Error()))
				return
			}
			return
		}

		var response Discord
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

func StartController(addr string) *controller.Controller {
	c, err := controller.NewController(addr)
	if err != nil {
		slog.Error("could not connect to the tor controller", slog.String("error", err.Error()))
		return nil
	}

	err = c.Authenticate("")
	if err != nil {
		slog.Error("could not authenticate with the tor controller", slog.String("error", err.Error()))
		return nil
	}

	slog.Info("connected to the tor controller")
	return c
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
