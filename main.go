package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const API = "https://api.discord.gx.games/v1/direct-fulfillment"
const Template = "https://discord.com/billing/partner-promotions/1180231712274387115/%s\n"

var client = http.DefaultClient
var file, _ = os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	go ListenForCircuit(30*time.Second, client)
	slog.Info("connected to the tor network")

	var discord Discord
	var resp *http.Response
	var err error

	Ticker(context.Background(), time.Second, func() bool {
		resp, err = client.Post(API, "application/json", bytes.NewBuffer([]byte("{\"partnerUserId\":\"cb7f04df-8b8e-4dc8-bc20-2b0e60e211d9\"}")))
		if err != nil {
			slog.Error("could not make request to discord api", slog.String("error", err.Error()))
			return false
		}

		err = json.NewDecoder(resp.Body).Decode(&discord)
		if err != nil {
			slog.Error("could not decode the discord api response", slog.String("error", err.Error()))
			return false
		}

		slog.Info("discord api response", slog.String("token", discord.Token))
		file.Write([]byte(fmt.Sprintf(Template, discord.Token)))

		return false
	})

	file.Close()
}

type Discord struct {
	Token string `json:"token"`
}
