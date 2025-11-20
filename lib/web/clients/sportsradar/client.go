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

func WithNbaAPIKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.nbaAPIKey = apiKey
	}
}

type Client struct {
	client *http.Client

	ncaafAPIKey  string
	ncaafBaseURL string

	ncaabAPIKey  string
	ncaabBaseURL string

	nbaAPIKey  string
	nbaBaseURL string
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		client:       &http.Client{},
		ncaafBaseURL: "https://api.sportradar.com/ncaafb/trial/v7",
		ncaabBaseURL: "https://api.sportradar.com/ncaamb/trial/v8",
		nbaBaseURL:   "https://api.sportradar.com/nba/trial/v8",
	}

	for _, option := range options {
		option(client)
	}

	return client
}
