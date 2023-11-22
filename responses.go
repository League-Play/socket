package main

type RequestUserInfoResponse struct {
	ResponseId string `json:"responseId"`
}

type FlowResponse struct {
	ResponseId string `json:"responseId"`
	Flow       string `json:"flow"`
}

type JoinLobbyResponse struct {
	ResponseId string `json:"responseId"`
	Username   string `json:"username"`
}

type ReadyResponse struct {
	ResponseId string `json:"responseId"`
	Username   string `json:"username"`
	IsReady    bool   `json:"isReady"`
}
