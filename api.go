package main

const API = "https://api.discord.gx.games/v1/direct-fulfillment"
const Template = "https://discord.com/billing/partner-promotions/1180231712274387115/%s\n"

type Discord struct {
	Token string `json:"token"`
}
