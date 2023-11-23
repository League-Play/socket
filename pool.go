package main

import (
	"fmt"
)

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
	Lobbys     []Lobby
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
	IsReady  bool
}

type Lobby struct {
	Users map[string]User
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Actions:    make(chan ClientAction),
		Users:      make([]User, 0),
		Lobbys:     []Lobby{{Users: make(map[string]User)}},
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
					// set userId as username (temporary for now)
					var user User = User{
						UserId:   uia.UserId,
						Flow:     "Home",
						Username: uia.UserId,
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
				// todo: integrate username with LeaguePlay app. Right now, it's just considering the userId to be the username
				var user *User = findUsername(pool.Users, jla.Username)
				if user != nil {
					ca.Client.Conn.WriteJSON(FlowResponse{ResponseId: "FlowResponse", Flow: "Lobby"})
					lobby := pool.Lobbys[0] // Assuming only one lobby for simplicity
					// Check if the user is already in the lobby
					if _, exists := lobby.Users[user.Username]; !exists {
						lobby.Users[user.Username] = *user // Add user to the lobby if not already present
					}
					// Send lobby response to the client
					for client, _ := range pool.Clients {
						for _, currentUser := range lobby.Users {
							fmt.Println(currentUser, currentUser.IsReady)
							// to do: add ready status to join lobby response for the case: user1 readies, user2 joins lobby but it shows user1 as not ready
							client.Conn.WriteJSON(JoinLobbyResponse{ResponseId: "JoinLobbyResponse", Username: currentUser.Username, IsReady: currentUser.IsReady})
						}
					}

				} else {
					// invalid username
				}
			case ReadyAction:
				var ra ReadyAction = a
				var user *User = findUsername(pool.Users, ra.Username)
				// change user status to ready in the lobby
				if user != nil {
					user.IsReady = ra.IsReady
					fmt.Println("user ready status: ", user.IsReady)
					pool.Lobbys[0].Users[user.Username] = *user
					for client, _ := range pool.Clients {
						client.Conn.WriteJSON(ReadyResponse{ResponseId: "ReadyResponse", Username: user.Username, IsReady: user.IsReady})
					}
				}
			}
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
