package main

import "fmt"

type ClientAction struct {
	Action Action
	Client *Client
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Actions    chan ClientAction
	Users      map[string]UserInfo
}

//Todo: refactor to group together into a Flow
// const (
// 	Home = iota
// 	Lobby = iota
// 	GameReport = iota
// )

// Todo: Refactor into enum
type UserInfo struct {
	Flow string
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Actions:    make(chan ClientAction),
		Users:      make(map[string]UserInfo),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	fmt.Println(client)
			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			// }
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			// fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			// }
			break
		case ca := <-pool.Actions:
			switch a := ca.Action.(type) {
			case UserInfoAction:
				var uia UserInfoAction = a
				if userInfo, exists := pool.Users[uia.UserId]; exists {
					fmt.Println("Writing response 1")
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: userInfo.Flow})
				} else {
					var userInfo UserInfo = UserInfo{
						Flow: "Home",
					}
					pool.Users[uia.UserId] = userInfo
					fmt.Println("Writing response")
					// for client, _ := range pool.Clients {
					// 	// client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
					// }
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: userInfo.Flow})
				}

			case RedirectAction:
				var ra RedirectAction = a
				if _, exists := pool.Users[ra.UserId]; exists {
					// Send back response
				} else {
					pool.Users[ra.UserId] = UserInfo{
						Flow: "Home",
					}
				}
			case JoinLobbyAction:

			}
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
