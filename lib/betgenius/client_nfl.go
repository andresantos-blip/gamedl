package betgenius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetNflSeasonsRaw(compId string) ([]byte, error) {
	url := fmt.Sprintf("%s/competitions/%s/seasons", c.fixturesV1URL, compId)
	return c.doV1Request(url)
}

func (c *Client) GetNflSeasons(compId string) (*SeasonsReply, error) {
	raw, err := c.GetNflSeasonsRaw(compId)
	if err != nil {
		return nil, fmt.Errorf("could not get nfl seasons: %w", err)
	}

	r := &SeasonsReply{}
	if err := json.Unmarshal(raw, r); err != nil {
		return nil, fmt.Errorf("could not unmarshal nfl seasons: %w", err)
	}

	return r, nil
}

func (c *Client) GetNflGamesForSeasonRaw(seasonID int) ([]byte, error) {
	url := fmt.Sprintf("%s/seasons/%d/fixtures", c.fixturesV1URL, seasonID)
	return c.doV1Request(url)
}

func (c *Client) GetNflGamesForSeason(seasonID int) (*GamesOfSeason, error) {
	raw, err := c.GetNflGamesForSeasonRaw(seasonID)
	if err != nil {
		return nil, fmt.Errorf("could not get nfl games for season: %w", err)
	}

	r := &GamesOfSeason{}
	if err := json.Unmarshal(raw, r); err != nil {
		return nil, fmt.Errorf("could not unmarshal nfl games for season: %w", err)
	}

	return r, nil
}

func (c *Client) GetNflPbpRaw(gameID string) ([]byte, error) {
	url := fmt.Sprintf("%s/fixtures/%s", c.fixturesV2URL, gameID)
	return c.doOAuthRequest(url)
}

func (c *Client) doV1Request(url string) ([]byte, error) {
	r, err := c.GetV1AuthedRequest(url)
	if err != nil {
		return nil, fmt.Errorf("could not get authed request: %w", err)
	}
	return c.doAndReadRequest(r)
}

func (c *Client) doOAuthRequest(url string) ([]byte, error) {
	r, err := c.GetOAuthAuthedRequest(url)
	if err != nil {
		return nil, fmt.Errorf("could not get authed request: %w", err)
	}
	return c.doAndReadRequest(r)
}

func (c *Client) doAndReadRequest(r *http.Request) ([]byte, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("could not get seasons: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body of seasons reply: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("bad status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
