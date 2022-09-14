package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type AuthenticationRequest struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func AuthenticateConnection(websocket **websocket.Conn, data json.RawMessage) bool {
	authRequest := AuthenticationRequest{}

	err := json.Unmarshal([]byte(data), &authRequest)
	if err != nil {
		fmt.Println("Error unmarshalling auth request:", err)
		return false
	}

	if authRequest.Name == "Browser Tester" && authRequest.APIKey == "1234567890" {
		return true
	}

	return false
}
