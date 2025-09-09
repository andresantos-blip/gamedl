package sportsradar

import (
	"net/http"
)

type ClientOption func(*Client)

func WithNcaafAPIKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.ncaafAPIKey = apiKey
	}
}

func WithNcaabAPIKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.ncaabAPIKey = apiKey
	}
}

type Client struct {
	client *http.Client

	ncaafAPIKey  string
	ncaafBaseURL string

	ncaabAPIKey  string
	ncaabBaseURL string
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		client:       &http.Client{},
		ncaafBaseURL: "https://api.sportradar.com/ncaafb/trial/v7",
		ncaabBaseURL: "https://api.sportradar.com/ncaamb/trial/v8",
	}

	for _, option := range options {
		option(client)
	}

	return client
}
