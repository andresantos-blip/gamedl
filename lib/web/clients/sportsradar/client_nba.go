package sportsradar

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetNbaSeasonsRaw() ([]byte, error) {
	url := fmt.Sprintf("%s/en/league/seasons.json?api_key=%s", c.nbaBaseURL, c.nbaAPIKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get seasons: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body of seasons reply: %w", err)
	}

	return body, nil
}

func (c *Client) GetNbaSeasons() (*NBASeasonsInfo, error) {
	body, err := c.GetNbaSeasonsRaw()
	if err != nil {
		return nil, err
	}
	seasons := &NBASeasonsInfo{}
	if err := json.Unmarshal(body, seasons); err != nil {
		return nil, fmt.Errorf("could not unmarshal seasons reply: %w", err)
	}
	return seasons, nil
}

func (c *Client) GetNbaSeasonScheduleRaw(year int) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%d/REG/schedule.json?api_key=%s", c.nbaBaseURL, year, c.nbaAPIKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get season schedule for year %d: %w", year, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body of season schedule reply: %w", err)
	}
	return body, nil
}

func (c *Client) GetNbaSeasonSchedule(year int) (*NbaSeasonSchedule, error) {
	body, err := c.GetNbaSeasonScheduleRaw(year)
	if err != nil {
		return nil, err
	}
	schedule := &NbaSeasonSchedule{}
	if err := json.Unmarshal(body, schedule); err != nil {
		return nil, fmt.Errorf("could not unmarshal season schedule reply: %w", err)
	}
	return schedule, nil
}

func (c *Client) GetNbaPbpOfGameRaw(gameId string) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%s/pbp.json?api_key=%s", c.nbaBaseURL, gameId, c.nbaAPIKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get play by play of game %s: %w", gameId, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body of play by play reply: %w", err)
	}
	return body, nil
}

func (c *Client) GetNbaPbpOfGame(gameId string) (*NbaGamePbp, error) {
	body, err := c.GetNbaPbpOfGameRaw(gameId)
	if err != nil {
		return nil, err
	}
	pbp := &NbaGamePbp{}
	if err := json.Unmarshal(body, pbp); err != nil {
		return nil, fmt.Errorf("could not unmarshal play by play reply: %w", err)
	}
	return pbp, nil
}
