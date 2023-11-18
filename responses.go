package main

type RequestUserInfoResponse struct {
	ResponseId string `json:"responseId"`
}

type FlowResponse struct {
	ResponseId string `json:"responseId"`
	Flow string `json:"flow"`
}
