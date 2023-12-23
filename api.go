package main

const API = "https://api.discord.gx.games/v1/direct-fulfillment"
const UserID = "cb7f04df-8b8e-4dc8-bc20-2b0e60e211d9"
const Template = "https://discord.com/billing/partner-promotions/1180231712274387115/%s\n"

type Discord struct {
	Token string `json:"token"`
}
