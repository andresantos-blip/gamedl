package sportsradar

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetNcaabSeasonsRaw() ([]byte, error) {
	url := fmt.Sprintf("%s/en/league/seasons.json?api_key=%s", c.ncaabBaseURL, c.ncaabAPIKey)
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

func (c *Client) GetNcaabSeasons() (*NcaabSeasonsInfo, error) {
	body, err := c.GetNcaabSeasonsRaw()
	if err != nil {
		return nil, err
	}
	seasons := &NcaabSeasonsInfo{}
	if err := json.Unmarshal(body, seasons); err != nil {
		return nil, fmt.Errorf("could not unmarshal seasons reply: %w", err)
	}
	return seasons, nil
}

func (c *Client) GetNcaabSeasonScheduleRaw(year int) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%d/REG/schedule.json?api_key=%s", c.ncaabBaseURL, year, c.ncaabAPIKey)
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

func (c *Client) GetNcaabSeasonSchedule(year int) (*NcaabSeasonSchedule, error) {
	body, err := c.GetNcaabSeasonScheduleRaw(year)
	if err != nil {
		return nil, err
	}
	schedule := &NcaabSeasonSchedule{}
	if err := json.Unmarshal(body, schedule); err != nil {
		return nil, fmt.Errorf("could not unmarshal season schedule reply: %w", err)
	}
	return schedule, nil
}

func (c *Client) GetNcaabPbpOfGameRaw(gameId string) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%s/pbp.json?api_key=%s", c.ncaabBaseURL, gameId, c.ncaabAPIKey)
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

func (c *Client) GetNcaabPbpOfGame(gameId string) (*NcaabGamePbp, error) {
	body, err := c.GetNcaabPbpOfGameRaw(gameId)
	if err != nil {
		return nil, err
	}
	pbp := &NcaabGamePbp{}
	if err := json.Unmarshal(body, pbp); err != nil {
		return nil, fmt.Errorf("could not unmarshal play by play reply: %w", err)
	}
	return pbp, nil
}
