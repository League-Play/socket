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
	Users      []User
}

//Todo: refactor to group together into a Flow
// const (
// 	Home = iota
// 	Lobby = iota
// 	GameReport = iota
// )

// Todo: Refactor into enum
type User struct {
	UserId   string
	Username string
	Flow     string
}

type Lobby struct {
	Users []User
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Actions:    make(chan ClientAction),
		Users:      make([]User, 0),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case ca := <-pool.Actions:
			switch a := ca.Action.(type) {
			case UserInfoAction:
				var uia UserInfoAction = a
				var user *User = findUser(pool.Users, uia.UserId)
				if user != nil {
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: user.Flow})
				} else {
					var user User = User{
						UserId:   uia.UserId,
						Flow:     "Home",
						Username: "username",
					}
					pool.Users = append(pool.Users, user)
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: user.Flow})
				}

			case RedirectAction:
				var ra RedirectAction = a
				var user *User = findUser(pool.Users, ra.UserId)
				if user != nil {
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: user.Flow})
				} else {
					var user User = User{
						UserId: ra.UserId,
						Flow:   "Home",
					}
					pool.Users = append(pool.Users, user)
				}
			case JoinLobbyAction:
				var jla JoinLobbyAction = a
				// todo: integrate username with LeaguePlay app. Right now, it's just expecting the hardcoded "username"
				var user *User = findUsername(pool.Users, jla.Username)
				if user != nil {
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: "Lobby"})
				} else {
					// invalid username
				}
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

func findUser(users []User, userId string) *User {
	for _, user := range users {
		if user.UserId == userId {
			return &user
		}
	}
	return nil
}

func findUsername(users []User, username string) *User {
	for _, user := range users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}
