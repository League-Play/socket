package main

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Actions    chan Action
	Users      map[string]UserInfo
}

type Flow int

// Todo: Refactor into enum
type UserInfo struct {
	Flow string
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Actions:    make(chan Action),
		Users:      make(map[string]UserInfo),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case action := <-pool.Actions:
			switch a := action.(type) {
			case RedirectAction:
				var ra RedirectAction = a
				if value, exists := pool.Users[ra.UserId]; exists {
					// Send back response
				} else {
					pool.Users[ra.UserId] = UserInfo{
						Flow: "Home",
					}
				}
			case JoinLobbyAction:

			}
			fmt.Printf("%T\n", action)
			// fmt.Println("Sending message to all clients in Pool")
			// for client, _ := range pool.Clients {
			//     if err := client.Conn.WriteJSON(message); err != nil {
			//         fmt.Println(err)
			//         return
			//     }
			// }
		}
	}
}
