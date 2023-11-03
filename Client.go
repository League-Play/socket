package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// Outbound Message
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Todo: DRY
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ur mom", string(p))
		var action ActionWrapper
		if err := json.Unmarshal(p, &action); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}
		fmt.Printf("Message Received: %+v\n", action)
		switch action.ActionId {
		case "Redirect":
			var redirectAction RedirectAction
			if err := json.Unmarshal(p, &redirectAction); err != nil {
				fmt.Println("Error decoding JSON: ", err)
				return
			}
			fmt.Printf("Message Received: %+v\n", redirectAction)
			c.Pool.Actions <- redirectAction
		case "JoinLobby":
			var joinLobbyAction JoinLobbyAction
			if err := json.Unmarshal(p, &joinLobbyAction); err != nil {
				fmt.Println("Error decoding JSON: ", err)
				return
			}
			fmt.Printf("Message Received: %+v\n", joinLobbyAction)
			c.Pool.Actions <- joinLobbyAction

		case "Ready":
			var readyAction ReadyAction
			if err := json.Unmarshal(p, &readyAction); err != nil {
				fmt.Println("Error decoding JSON: ", err)
				return
			}
			fmt.Printf("Message Received: %+v\n", readyAction)
			c.Pool.Actions <- readyAction
		}
	}
}
