package sportsradar

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetNcaafSeasonsRaw() ([]byte, error) {
	url := fmt.Sprintf("%s/en/league/seasons.json?api_key=%s", c.ncaafBaseURL, c.ncaafAPIKey)
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

func (c *Client) GetNcaafSeasons() (*NcaafSeasonsInfo, error) {
	body, err := c.GetNcaafSeasonsRaw()
	if err != nil {
		return nil, err
	}
	seasons := &NcaafSeasonsInfo{}
	if err := json.Unmarshal(body, seasons); err != nil {
		return nil, fmt.Errorf("could not unmarshal seasons reply: %w", err)
	}
	return seasons, nil
}

func (c *Client) GetNcaafSeasonScheduleRaw(year int) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%d/REG/schedule.json?api_key=%s", c.ncaafBaseURL, year, c.ncaafAPIKey)
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

func (c *Client) GetNcaafSeasonSchedule(year int) (*NcaafSeasonSchedule, error) {
	body, err := c.GetNcaafSeasonScheduleRaw(year)
	if err != nil {
		return nil, err
	}
	schedule := &NcaafSeasonSchedule{}
	if err := json.Unmarshal(body, schedule); err != nil {
		return nil, fmt.Errorf("could not unmarshal season schedule reply: %w", err)
	}
	return schedule, nil
}

func (c *Client) GetNcaafPbpOfGameRaw(gameId string) ([]byte, error) {
	url := fmt.Sprintf("%s/en/games/%s/pbp.json?api_key=%s", c.ncaafBaseURL, gameId, c.ncaafAPIKey)
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

func (c *Client) GetNcaafPbpOfGame(gameId string) (*NcaafGamePbp, error) {
	body, err := c.GetNcaafPbpOfGameRaw(gameId)
	if err != nil {
		return nil, err
	}
	pbp := &NcaafGamePbp{}
	if err := json.Unmarshal(body, pbp); err != nil {
		return nil, fmt.Errorf("could not unmarshal play by play reply: %w", err)
	}
	return pbp, nil
}
