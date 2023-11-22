package main

type Action interface {
	Type() string
}

type ActionWrapper struct {
	ActionId string `json:"action"`
}

func (aw ActionWrapper) Type() string {
	return aw.ActionId
}

type RedirectAction struct {
	ActionId string `json:"action"`
	UserId   string `json:"user"`
}

func (ra RedirectAction) Type() string {
	return ra.ActionId
}

type JoinLobbyAction struct {
	ActionId string `json:"action"`
	Username string `json:"username"`
}

func (jla JoinLobbyAction) Type() string {
	return jla.ActionId
}

type ReadyAction struct {
	ActionId string `json:"action"`
	Username string `json:"username"`
	IsReady  bool   `json:"isReady"`
}

func (ra ReadyAction) Type() string {
	return ra.ActionId
}

type UserInfoAction struct {
	UserId string `json:"user"`
}

func (uia UserInfoAction) Type() string {
	return uia.UserId
}
