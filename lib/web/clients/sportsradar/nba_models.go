package sportsradar

import (
	"slices"
	"time"
)

type NBASeasonsInfo struct {
	League struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"league"`
	Seasons []*NBASeasonInfo `json:"seasons"`
}

func (si *NBASeasonsInfo) Years() []int {
	years := make([]int, 0, len(si.Seasons))
	for _, season := range si.Seasons {
		years = append(years, season.Year)
	}
	return years
}

func (si *NBASeasonsInfo) FilterYears(years []int) {
	filteredSeasons := make([]*NBASeasonInfo, 0, len(si.Seasons))

	for _, season := range si.Seasons {
		if slices.Contains(years, season.Year) {
			filteredSeasons = append(filteredSeasons, season)
		}
	}
	si.Seasons = filteredSeasons
}

func (si *NBASeasonsInfo) FilterSeasonType(seasonType string) {
	filteredSeasons := make([]*NBASeasonInfo, 0, len(si.Seasons))

	for _, season := range si.Seasons {
		if season.Type.Code == seasonType {
			filteredSeasons = append(filteredSeasons, season)
		}
	}
	si.Seasons = filteredSeasons
}

type NBASeasonInfo struct {
	Id   string `json:"id"`
	Year int    `json:"year"`
	Type struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Status    string `json:"status,omitempty"`
}

type NbaSeasonSchedule struct {
	League struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"league"`
	Season struct {
		Id   string `json:"id"`
		Year int    `json:"year"`
		Type string `json:"type"`
	} `json:"season"`
	Games []*NBAGame `json:"games"`
}

type NBAGame struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	Coverage     string    `json:"coverage"`
	Scheduled    time.Time `json:"scheduled"`
	HomePoints   int       `json:"home_points,omitempty"`
	AwayPoints   int       `json:"away_points,omitempty"`
	TrackOnCourt bool      `json:"track_on_court"`
	SrId         string    `json:"sr_id,omitempty"`
	Reference    string    `json:"reference,omitempty"`
	TimeZones    struct {
		Venue string `json:"venue"`
		Home  string `json:"home,omitempty"`
		Away  string `json:"away,omitempty"`
	} `json:"time_zones,omitempty"`
	Venue struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
		Address  string `json:"address"`
		City     string `json:"city"`
		State    string `json:"state,omitempty"`
		Zip      string `json:"zip,omitempty"`
		Country  string `json:"country"`
		SrId     string `json:"sr_id"`
		Location struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"location"`
	} `json:"venue,omitempty"`
	Broadcasts []struct {
		Network string `json:"network"`
		Type    string `json:"type"`
		Locale  string `json:"locale,omitempty"`
		Channel string `json:"channel,omitempty"`
	} `json:"broadcasts,omitempty"`
	Home struct {
		Name      string `json:"name"`
		Alias     string `json:"alias"`
		Id        string `json:"id"`
		SrId      string `json:"sr_id,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"home"`
	Away struct {
		Name      string `json:"name"`
		Alias     string `json:"alias"`
		Id        string `json:"id"`
		SrId      string `json:"sr_id,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"away"`
	InseasonTournament bool   `json:"inseason_tournament,omitempty"`
	NeutralSite        bool   `json:"neutral_site,omitempty"`
	Title              string `json:"title,omitempty"`
}

type NbaGamePbp struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	Coverage     string    `json:"coverage"`
	Scheduled    time.Time `json:"scheduled"`
	Duration     string    `json:"duration"`
	Attendance   int       `json:"attendance"`
	LeadChanges  int       `json:"lead_changes"`
	TimesTied    int       `json:"times_tied"`
	Clock        string    `json:"clock"`
	Quarter      int       `json:"quarter"`
	TrackOnCourt bool      `json:"track_on_court"`
	Reference    string    `json:"reference"`
	EntryMode    string    `json:"entry_mode"`
	SrID         string    `json:"sr_id"`
	ClockDecimal string    `json:"clock_decimal"`
	Broadcasts   []struct {
		Type    string `json:"type"`
		Locale  string `json:"locale"`
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
		Bonus             bool   `json:"bonus"`
		SrID              string `json:"sr_id"`
		RemainingTimeouts int    `json:"remaining_timeouts"`
		Reference         string `json:"reference"`
		Record            struct {
			Wins   int `json:"wins"`
			Losses int `json:"losses"`
		} `json:"record"`
	} `json:"home"`
	Away struct {
		Name              string `json:"name"`
		Alias             string `json:"alias"`
		Market            string `json:"market"`
		ID                string `json:"id"`
		Points            int    `json:"points"`
		Bonus             bool   `json:"bonus"`
		SrID              string `json:"sr_id"`
		RemainingTimeouts int    `json:"remaining_timeouts"`
		Reference         string `json:"reference"`
		Record            struct {
			Wins   int `json:"wins"`
			Losses int `json:"losses"`
		} `json:"record"`
	} `json:"away"`
	Periods []struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Number   int    `json:"number"`
		Sequence int    `json:"sequence"`
		Scoring  struct {
			LeadChanges int `json:"lead_changes"`
			Home        struct {
				Name      string `json:"name"`
				Market    string `json:"market"`
				ID        string `json:"id"`
				Points    int    `json:"points"`
				Reference string `json:"reference"`
			} `json:"home"`
			Away struct {
				Name      string `json:"name"`
				Market    string `json:"market"`
				ID        string `json:"id"`
				Points    int    `json:"points"`
				Reference string `json:"reference"`
			} `json:"away"`
		} `json:"scoring"`
		Events []struct {
			ID           string    `json:"id"`
			Clock        string    `json:"clock"`
			Updated      time.Time `json:"updated"`
			Description  string    `json:"description"`
			WallClock    time.Time `json:"wall_clock"`
			Sequence     int64     `json:"sequence"`
			HomePoints   int       `json:"home_points"`
			AwayPoints   int       `json:"away_points"`
			ClockDecimal string    `json:"clock_decimal"`
			Created      time.Time `json:"created"`
			Number       int       `json:"number"`
			EventType    string    `json:"event_type"`
			Attribution  struct {
				Name      string `json:"name"`
				Market    string `json:"market"`
				ID        string `json:"id"`
				SrID      string `json:"sr_id"`
				Reference string `json:"reference"`
			} `json:"attribution,omitempty"`
			OnCourt struct {
				Home struct {
					Name      string `json:"name"`
					Market    string `json:"market"`
					ID        string `json:"id"`
					SrID      string `json:"sr_id"`
					Reference string `json:"reference"`
					Players   []struct {
						FullName     string `json:"full_name"`
						JerseyNumber string `json:"jersey_number"`
						ID           string `json:"id"`
						SrID         string `json:"sr_id"`
						Reference    string `json:"reference"`
					} `json:"players"`
				} `json:"home"`
				Away struct {
					Name      string `json:"name"`
					Market    string `json:"market"`
					ID        string `json:"id"`
					SrID      string `json:"sr_id"`
					Reference string `json:"reference"`
					Players   []struct {
						FullName     string `json:"full_name"`
						JerseyNumber string `json:"jersey_number"`
						ID           string `json:"id"`
						SrID         string `json:"sr_id"`
						Reference    string `json:"reference"`
					} `json:"players"`
				} `json:"away"`
			} `json:"on_court,omitempty"`
			Possession struct {
				Name      string `json:"name"`
				Market    string `json:"market"`
				ID        string `json:"id"`
				SrID      string `json:"sr_id"`
				Reference string `json:"reference"`
			} `json:"possession,omitempty"`
			Location struct {
				CoordX     int    `json:"coord_x"`
				CoordY     int    `json:"coord_y"`
				ActionArea string `json:"action_area"`
			} `json:"location,omitempty"`
			Statistics []struct {
				Type           string  `json:"type"`
				Made           bool    `json:"made"`
				ShotType       string  `json:"shot_type"`
				ThreePointShot bool    `json:"three_point_shot"`
				ShotDistance   float64 `json:"shot_distance"`
				Team           *struct {
					Name      string `json:"name"`
					Market    string `json:"market"`
					ID        string `json:"id"`
					SrID      string `json:"sr_id"`
					Reference string `json:"reference"`
				} `json:"team,omitempty"`
				Player *struct {
					FullName     string `json:"full_name"`
					JerseyNumber string `json:"jersey_number"`
					ID           string `json:"id"`
					SrID         string `json:"sr_id"`
					Reference    string `json:"reference"`
				} `json:"player,omitempty"`
			} `json:"statistics,omitempty"`
			Qualifiers []struct {
				Qualifier string `json:"qualifier"`
			} `json:"qualifiers,omitempty"`
			Attempt      string `json:"attempt,omitempty"`
			TurnoverType string `json:"turnover_type,omitempty"`
			Duration     int    `json:"duration,omitempty"`
		} `json:"events"`
	} `json:"periods"`
	DeletedEvents []struct {
		ID string `json:"id"`
	} `json:"deleted_events"`
}
