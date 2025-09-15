package sportsradar

import (
	"slices"
	"time"
)

type NcaabSeasonInfo struct {
	ID   string `json:"id"`
	Year int    `json:"year"`
	Type struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Status    string `json:"status,omitempty"`
}
type NcaabSeasonsInfo struct {
	League struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"league"`
	Seasons []*NcaabSeasonInfo `json:"seasons"`
}

func (si *NcaabSeasonsInfo) Years() []int {
	years := make([]int, 0, len(si.Seasons))
	for _, season := range si.Seasons {
		years = append(years, season.Year)
	}
	return years
}

func (si *NcaabSeasonsInfo) FilterYears(years []int) {
	filteredSeasons := make([]*NcaabSeasonInfo, 0, len(si.Seasons))

	for _, season := range si.Seasons {
		if slices.Contains(years, season.Year) {
			filteredSeasons = append(filteredSeasons, season)
		}
	}
	si.Seasons = filteredSeasons
}

type NcaabSeasonSchedule struct {
	League struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"league"`
	Season struct {
		ID   string `json:"id"`
		Year int    `json:"year"`
		Type string `json:"type"`
	} `json:"season"`
	Games []*NcaabGame `json:"games"`
}

type NcaabGame struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"`
	Coverage       string    `json:"coverage"`
	Scheduled      time.Time `json:"scheduled"`
	HomePoints     int       `json:"home_points,omitempty"`
	AwayPoints     int       `json:"away_points,omitempty"`
	ConferenceGame bool      `json:"conference_game"`
	TimeZones      struct {
		Venue string `json:"venue"`
		Home  string `json:"home"`
	} `json:"time_zones"`
	Venue struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
		Address  string `json:"address"`
		City     string `json:"city"`
		State    string `json:"state"`
		Zip      string `json:"zip"`
		Country  string `json:"country"`
		Location struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"location"`
	} `json:"venue"`
	Broadcasts []struct {
		Network string `json:"network"`
		Type    string `json:"type"`
	} `json:"broadcasts,omitempty"`
	Home struct {
		Name  string `json:"name"`
		Alias string `json:"alias"`
		ID    string `json:"id"`
	} `json:"home"`
	Away struct {
		Name  string `json:"name"`
		Alias string `json:"alias"`
		ID    string `json:"id"`
	} `json:"away"`
	NeutralSite  bool   `json:"neutral_site,omitempty"`
	TrackOnCourt bool   `json:"track_on_court,omitempty"`
	Title        string `json:"title,omitempty"`
}

type NcaabGamePbp struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"`
	Coverage       string    `json:"coverage"`
	NeutralSite    bool      `json:"neutral_site"`
	Scheduled      time.Time `json:"scheduled"`
	ConferenceGame bool      `json:"conference_game"`
	Attendance     int       `json:"attendance"`
	LeadChanges    int       `json:"lead_changes"`
	TimesTied      int       `json:"times_tied"`
	Clock          string    `json:"clock"`
	Half           int       `json:"half"`
	TrackOnCourt   bool      `json:"track_on_court"`
	EntryMode      string    `json:"entry_mode"`
	ClockDecimal   string    `json:"clock_decimal"`
	Broadcasts     []struct {
		Type    string `json:"type"`
		Network string `json:"network"`
	} `json:"broadcasts"`
	TimeZones struct {
		Venue string `json:"venue"`
		Home  string `json:"home"`
		Away  string `json:"away"`
	} `json:"time_zones"`
	Season struct {
		ID   string `json:"id"`
		Year int    `json:"year"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"season"`
	Home struct {
		Name              string `json:"name"`
		Alias             string `json:"alias"`
		Market            string `json:"market"`
		ID                string `json:"id"`
		Points            int    `json:"points"`
		Rank              int    `json:"rank"`
		RemainingTimeouts int    `json:"remaining_timeouts"`
	} `json:"home"`
	Away struct {
		Name              string `json:"name"`
		Alias             string `json:"alias"`
		Market            string `json:"market"`
		ID                string `json:"id"`
		Points            int    `json:"points"`
		Rank              int    `json:"rank"`
		RemainingTimeouts int    `json:"remaining_timeouts"`
	} `json:"away"`
	Periods []struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Number   int    `json:"number"`
		Sequence int    `json:"sequence"`
		Scoring  struct {
			TimesTied   int `json:"times_tied"`
			LeadChanges int `json:"lead_changes"`
			Home        struct {
				Name   string `json:"name"`
				Market string `json:"market"`
				ID     string `json:"id"`
				Points int    `json:"points"`
			} `json:"home"`
			Away struct {
				Name   string `json:"name"`
				Market string `json:"market"`
				ID     string `json:"id"`
				Points int    `json:"points"`
			} `json:"away"`
		} `json:"scoring"`
		Events []struct {
			ID           string    `json:"id"`
			Clock        string    `json:"clock"`
			Updated      time.Time `json:"updated"`
			Description  string    `json:"description"`
			Sequence     int64     `json:"sequence"`
			HomePoints   int       `json:"home_points"`
			AwayPoints   int       `json:"away_points"`
			ClockDecimal string    `json:"clock_decimal"`
			Created      time.Time `json:"created"`
			EventType    string    `json:"event_type"`
			Attribution  struct {
				Name       string `json:"name"`
				Market     string `json:"market"`
				ID         string `json:"id"`
				TeamBasket string `json:"team_basket"`
			} `json:"attribution,omitempty"`
			Location struct {
				CoordX int `json:"coord_x"`
				CoordY int `json:"coord_y"`
			} `json:"location,omitempty"`
			Possession struct {
				Name   string `json:"name"`
				Market string `json:"market"`
				ID     string `json:"id"`
			} `json:"possession,omitempty"`
			Statistics []struct {
				Type string `json:"type"`
				Team struct {
					Name   string `json:"name"`
					Market string `json:"market"`
					ID     string `json:"id"`
				} `json:"team"`
				Player struct {
					FullName     string `json:"full_name"`
					JerseyNumber string `json:"jersey_number"`
					ID           string `json:"id"`
				} `json:"player"`
			} `json:"statistics,omitempty"`
			TurnoverType string `json:"turnover_type,omitempty"`
			Attempt      string `json:"attempt,omitempty"`
			Duration     int    `json:"duration,omitempty"`
		} `json:"events"`
	} `json:"periods"`
	DeletedEvents []struct {
		ID string `json:"id"`
	} `json:"deleted_events"`
}
