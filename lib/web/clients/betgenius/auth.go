package betgenius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type AuthV1Reply struct {
	AccessToken  string `json:"AccessToken"`
	ExpiresIn    int    `json:"ExpiresIn"`
	TokenType    string `json:"TokenType"`
	RefreshToken string `json:"RefreshToken"`
	IDToken      string `json:"IdToken"`
}

type UsernamePassword struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type OAuthReply struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type OAuthToken struct {
	m         *sync.RWMutex
	Raw       OAuthReply
	ExpiresIn *time.Time
}

func NewOAuthToken() *OAuthToken {
	return &OAuthToken{
		m: &sync.RWMutex{},
	}
}

func (t *OAuthToken) IsExpired() bool {
	t.m.RLock()
	defer t.m.RUnlock()
	if t.ExpiresIn == nil {
		return true
	}
	return time.Now().After(*t.ExpiresIn)
}

type TokenV1 struct {
	m         *sync.RWMutex
	Raw       AuthV1Reply
	ExpiresIn *time.Time
}

func NewTokenV1() *TokenV1 {
	return &TokenV1{
		m: &sync.RWMutex{},
	}
}

func (t *TokenV1) IsExpired() bool {
	t.m.RLock()
	defer t.m.RUnlock()
	if t.ExpiresIn == nil {
		return true
	}
	return time.Now().After(*t.ExpiresIn)
}

func (c *Client) GetV1Token() (string, error) {
	c.v1Token.m.RLock()
	if c.v1Token != nil && !c.v1Token.IsExpired() {
		c.v1Token.m.RUnlock()
		return c.v1Token.Raw.IDToken, nil
	}
	c.v1Token.m.RUnlock()

	auth := &UsernamePassword{
		Username: c.fixtureUsername,
		Password: c.fixturePassword,
	}

	payload, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal auth: %w", err)
	}

	buff := bytes.NewBuffer(payload)

	req, err := http.NewRequest("POST", c.authV1, buff)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "PostmanRuntime/7.26.8")

	r, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	startExp := time.Now()
	reply := &AuthV1Reply{}
	replyData, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("could not read auth v1 reply: %w", err)
	}
	if err := json.Unmarshal(replyData, reply); err != nil {
		return "", fmt.Errorf("could not unmarshal auth v1 reply: %w", err)
	}

	c.v1Token.m.Lock()
	defer c.v1Token.m.Unlock()

	now := time.Now()
	elapsed := now.Sub(startExp)
	expiresIn := now.Add(time.Duration(reply.ExpiresIn)*time.Second - elapsed)
	c.v1Token.ExpiresIn = &expiresIn
	c.v1Token.Raw = *reply

	return reply.IDToken, nil
}

func (c *Client) GetOAuthToken() (string, error) {
	c.oAuthToken.m.RLock()
	if c.oAuthToken != nil && !c.oAuthToken.IsExpired() {
		c.oAuthToken.m.RUnlock()
		return c.oAuthToken.Raw.AccessToken, nil
	}
	c.oAuthToken.m.RUnlock()

	req, err := http.NewRequest("POST", c.authOauth, nil)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	req.SetBasicAuth(c.statsUsername, c.statsPassword)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	startExp := time.Now()
	reply := &OAuthReply{}

	replyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read auth v1 reply: %w", err)
	}

	if err = json.Unmarshal(replyData, reply); err != nil {
		return "", fmt.Errorf("could not unmarshal auth v1 reply: %w", err)
	}

	c.oAuthToken.m.Lock()
	defer c.oAuthToken.m.Unlock()

	now := time.Now()
	elapsed := now.Sub(startExp)
	expiresIn := now.Add(time.Duration(reply.ExpiresIn)*time.Second - elapsed)
	c.oAuthToken.ExpiresIn = &expiresIn
	c.oAuthToken.Raw = *reply

	return reply.AccessToken, nil
}

func (c *Client) GetV1AuthedRequest(url string) (*http.Request, error) {
	token, err := c.GetV1Token()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Api-Key", c.fixtureKey)

	return req, nil
}

func (c *Client) GetOAuthAuthedRequest(url string) (*http.Request, error) {
	token, err := c.GetOAuthToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Api-Key", c.statsKey)

	return req, nil
}
