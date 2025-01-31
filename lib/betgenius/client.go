package betgenius

import "net/http"

type ClientOption func(*Client)

func WithFixtureKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.fixtureKey = apiKey
	}
}

func WithFixtureUsername(username string) ClientOption {
	return func(client *Client) {
		client.fixtureUsername = username
	}
}

func WithFixturePassword(password string) ClientOption {
	return func(client *Client) {
		client.fixturePassword = password
	}
}

func WithStatsKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.statsKey = apiKey
	}
}

func WithStatsUsername(username string) ClientOption {
	return func(client *Client) {
		client.statsUsername = username
	}
}

func WithStatsPassword(password string) ClientOption {
	return func(client *Client) {
		client.statsPassword = password
	}
}

type Client struct {
	client *http.Client

	fixtureKey      string
	fixtureUsername string
	fixturePassword string

	statsKey      string
	statsUsername string
	statsPassword string

	authV1    string
	authOauth string

	fixturesV1URL string
	fixturesV2URL string

	v1Token    *TokenV1
	oAuthToken *OAuthToken
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		client:        &http.Client{},
		authV1:        "https://api.geniussports.com/Auth-v1/PROD/login",
		authOauth:     "https://auth.api.geniussports.com/oauth2/token?grant_type=client_credentials&scope=statistics-api%2Fstatistics%3Aread%20statistics-api%2Fliveaccess%3Aread%20matchstateapi%2Fmatchstate%3Aread%20matchstateapi%2Fgranularity%3Aread",
		fixturesV1URL: "https://api.geniussports.com/Fixtures-v1/PRODPRM",
		fixturesV2URL: "https://platform.matchstate.api.geniussports.com/api/v2/sources/GeniusPremium/sports/17",
		v1Token:       NewTokenV1(),
		oAuthToken:    NewOAuthToken(),
	}

	for _, option := range options {
		option(client)
	}

	return client
}
