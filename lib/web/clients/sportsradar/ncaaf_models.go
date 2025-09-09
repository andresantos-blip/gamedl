package sportsradar

import "time"

type NcaafSeasonsInfo struct {
	League struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"league"`
	Seasons []struct {
		ID        string `json:"id"`
		Year      int    `json:"year"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
		Type      struct {
			Code string `json:"code"`
		} `json:"type"`
	} `json:"seasons"`
}

func (si *NcaafSeasonsInfo) Years() []int {
	years := make([]int, 0, len(si.Seasons))
	for _, season := range si.Seasons {
		years = append(years, season.Year)
	}
	return years
}

type NcaafSeasonSchedule struct {
	ID    string `json:"id"`
	Year  int    `json:"year"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Weeks []struct {
		ID       string       `json:"id"`
		Sequence int          `json:"sequence"`
		Title    string       `json:"title"`
		Games    []*NcaafGame `json:"games"`
		ByeWeek  []struct {
			Team struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Alias string `json:"alias"`
			} `json:"team"`
		} `json:"bye_week"`
	} `json:"weeks"`
	Comment string `json:"_comment"`
}

type NcaafGame struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"`
	Scheduled      time.Time `json:"scheduled"`
	Attendance     int       `json:"attendance,omitempty"`
	EntryMode      string    `json:"entry_mode"`
	Coverage       string    `json:"coverage"`
	NeutralSite    bool      `json:"neutral_site,omitempty"`
	GameType       string    `json:"game_type"`
	ConferenceGame bool      `json:"conference_game"`
	Duration       string    `json:"duration,omitempty"`
	Venue          struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Address  string `json:"address"`
		Capacity int    `json:"capacity"`
		Surface  string `json:"surface"`
		RoofType string `json:"roof_type"`
		Location struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"location"`
	} `json:"venue"`
	Home struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Alias      string `json:"alias"`
		GameNumber int    `json:"game_number"`
	} `json:"home"`
	Away struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Alias      string `json:"alias"`
		GameNumber int    `json:"game_number"`
	} `json:"away"`
	Broadcast struct {
		Network string `json:"network"`
	} `json:"broadcast"`
	TimeZones struct {
		Venue string `json:"venue"`
		Home  string `json:"home"`
		Away  string `json:"away"`
	} `json:"time_zones"`
	Weather struct {
		Condition string `json:"condition"`
		Humidity  int    `json:"humidity"`
		Temp      int    `json:"temp"`
		Wind      struct {
			Speed     int    `json:"speed"`
			Direction string `json:"direction"`
		} `json:"wind"`
	} `json:"weather"`
	Scoring struct {
		HomePoints int `json:"home_points"`
		AwayPoints int `json:"away_points"`
		Periods    []struct {
			PeriodType string `json:"period_type"`
			ID         string `json:"id"`
			Number     int    `json:"number"`
			Sequence   int    `json:"sequence"`
			HomePoints int    `json:"home_points"`
			AwayPoints int    `json:"away_points"`
		} `json:"periods"`
	} `json:"scoring"`
}

type NcaafGamePbp struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"`
	Scheduled      time.Time `json:"scheduled"`
	Attendance     int       `json:"attendance"`
	EntryMode      string    `json:"entry_mode"`
	Clock          string    `json:"clock"`
	Quarter        int       `json:"quarter"`
	Coverage       string    `json:"coverage"`
	NeutralSite    bool      `json:"neutral_site"`
	GameType       string    `json:"game_type"`
	ConferenceGame bool      `json:"conference_game"`
	Title          string    `json:"title"`
	Duration       string    `json:"duration"`
	ParentID       string    `json:"parent_id"`
	Weather        *struct {
		Condition string `json:"condition"`
		Humidity  int    `json:"humidity"`
		Temp      int    `json:"temp"`
		Wind      *struct {
			Speed     int    `json:"speed"`
			Direction string `json:"direction"`
		} `json:"window,omitempty"`
	} `json:"weather,omitempty"`
	Summary *struct {
		Season *struct {
			ID   string `json:"id"`
			Year int    `json:"year"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"season,omitempty"`
		Week *struct {
			ID       string `json:"id"`
			Sequence int    `json:"sequence"`
			Title    string `json:"title"`
		} `json:"week,omitempty"`
		Venue *struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			City     string `json:"city"`
			State    string `json:"state"`
			Country  string `json:"country"`
			Zip      string `json:"zip"`
			Address  string `json:"address"`
			Capacity int    `json:"capacity"`
			Surface  string `json:"surface"`
			RoofType string `json:"roof_type"`
			SrID     string `json:"sr_id"`
			Location *struct {
				Lat string `json:"lat"`
				Lng string `json:"lng"`
			} `json:"location,omitempty"`
		} `json:"venue,omitempty"`
		Home *struct {
			ID                  string `json:"id"`
			Name                string `json:"name"`
			Market              string `json:"market"`
			Alias               string `json:"alias"`
			UsedTimeouts        int    `json:"used_timeouts"`
			RemainingTimeouts   int    `json:"remaining_timeouts"`
			Points              int    `json:"points"`
			UsedChallenges      int    `json:"used_challenges"`
			RemainingChallenges int    `json:"remaining_challenges"`
			Record              *struct {
				Wins   int `json:"wins"`
				Losses int `json:"losses"`
				Ties   int `json:"ties"`
			} `json:"record,omitempty"`
		} `json:"home,omitempty"`
		Away *struct {
			ID                  string `json:"id"`
			Name                string `json:"name"`
			Market              string `json:"market"`
			Alias               string `json:"alias"`
			UsedTimeouts        int    `json:"used_timeouts"`
			RemainingTimeouts   int    `json:"remaining_timeouts"`
			Points              int    `json:"points"`
			UsedChallenges      int    `json:"used_challenges"`
			RemainingChallenges int    `json:"remaining_challenges"`
			Record              *struct {
				Wins   int `json:"wins"`
				Losses int `json:"losses"`
				Ties   int `json:"ties"`
			} `json:"record,omitempty"`
		} `json:"away,omitempty"`
	} `json:"summary,omitempty"`
	Broadcast *struct {
		Network   string `json:"network"`
		Satellite string `json:"satellite"`
	} `json:"broadcast,omitempty"`
	Periods []struct {
		PeriodType string  `json:"period_type"`
		ID         string  `json:"id"`
		Number     int     `json:"number"`
		Sequence   float64 `json:"sequence"`
		Scoring    *struct {
			Home *struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Market string `json:"market"`
				Alias  string `json:"alias"`
				Points int    `json:"points"`
			} `json:"home,omitempty"`
			Away *struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Market string `json:"market"`
				Alias  string `json:"alias"`
				Points int    `json:"points"`
			} `json:"away,omitempty"`
		} `json:"scoring,omitempty"`
		CoinToss *struct {
			Home *struct {
				Outcome  string `json:"outcome"`
				Decision string `json:"decision"`
			} `json:"home,omitempty"`
			Away *struct {
				Outcome  string `json:"outcome"`
				Decision string `json:"decision"`
			} `json:"away,omitempty"`
		} `json:"coin_toss,omitempty"`
		Pbp []struct {
			Type                  string    `json:"type"`
			ID                    string    `json:"id"`
			Sequence              float64   `json:"sequence"`
			StartReason           string    `json:"start_reason"`
			EndReason             string    `json:"end_reason"`
			PlayCount             int       `json:"play_count"`
			Duration              string    `json:"duration"`
			FirstDowns            int       `json:"first_downs"`
			Gain                  int       `json:"gain"`
			PenaltyYards          int       `json:"penalty_yards"`
			CreatedAt             time.Time `json:"created_at"`
			UpdatedAt             time.Time `json:"updated_at"`
			TeamSequence          int       `json:"team_sequence"`
			StartClock            string    `json:"start_clock"`
			EndClock              string    `json:"end_clock"`
			FirstDriveYardline    int       `json:"first_drive_yardline"`
			LastDriveYardline     int       `json:"last_drive_yardline"`
			FarthestDriveYardline int       `json:"farthest_drive_yardline"`
			NetYards              int       `json:"net_yards"`
			PatPointsAttempted    int       `json:"pat_points_attempted"`
			OffensiveTeam         *struct {
				Points int    `json:"points"`
				ID     string `json:"id"`
			} `json:"offensive_team,omitempty"`
			DefensiveTeam *struct {
				Points int    `json:"points"`
				ID     string `json:"id"`
			} `json:"defensive_team,omitempty"`
			Events []struct {
				Type           string    `json:"type"`
				ID             string    `json:"id"`
				Sequence       float64   `json:"sequence"`
				Clock          string    `json:"clock"`
				HomePoints     int       `json:"home_points,omitempty"`
				AwayPoints     int       `json:"away_points,omitempty"`
				PlayType       string    `json:"play_type,omitempty"`
				WallClock      time.Time `json:"wall_clock"`
				Description    string    `json:"description"`
				FakePunt       bool      `json:"fake_punt,omitempty"`
				FakeFieldGoal  bool      `json:"fake_field_goal,omitempty"`
				ScreenPass     bool      `json:"screen_pass,omitempty"`
				PlayAction     bool      `json:"play_action,omitempty"`
				RunPassOption  bool      `json:"run_pass_option,omitempty"`
				CreatedAt      time.Time `json:"created_at"`
				UpdatedAt      time.Time `json:"updated_at"`
				StartSituation *struct {
					Clock      string `json:"clock"`
					Down       int    `json:"down"`
					Yfd        int    `json:"yfd"`
					Possession *struct {
						ID     string `json:"id"`
						Name   string `json:"name"`
						Market string `json:"market"`
						Alias  string `json:"alias"`
					} `json:"possession,omitempty"`
					Location *struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						Market   string `json:"market"`
						Alias    string `json:"alias"`
						Yardline int    `json:"yardline"`
					} `json:"location,omitempty"`
				} `json:"start_situation,omitempty"`
				EndSituation *struct {
					Clock      string `json:"clock"`
					Down       int    `json:"down"`
					Yfd        int    `json:"yfd"`
					Possession *struct {
						ID     string `json:"id"`
						Name   string `json:"name"`
						Market string `json:"market"`
						Alias  string `json:"alias"`
					} `json:"possession,omitempty"`
					Location *struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						Market   string `json:"market"`
						Alias    string `json:"alias"`
						Yardline int    `json:"yardline"`
					} `json:"location,omitempty"`
				} `json:"end_situation,omitempty"`
				Statistics []struct {
					StatType string `json:"stat_type"`
					Attempt  int    `json:"attempt,omitempty"`
					Yards    int    `json:"yards,omitempty"`
					NetYards int    `json:"net_yards,omitempty"`
					Endzone  int    `json:"endzone,omitempty"`
					Player   *struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						Jersey   string `json:"jersey"`
						Position string `json:"position"`
					} `json:"player,omitempty"`
					Team *struct {
						ID     string `json:"id"`
						Name   string `json:"name"`
						Market string `json:"market"`
						Alias  string `json:"alias"`
					} `json:"team,omitempty"`
					Return   int    `json:"return,omitempty"`
					Category string `json:"category,omitempty"`
					Tackle   int    `json:"tackle,omitempty"`
				} `json:"statistics,omitempty"`
				Details []struct {
					Category      string  `json:"category"`
					Description   string  `json:"description"`
					Sequence      float64 `json:"sequence"`
					Yards         int     `json:"yards,omitempty"`
					Result        string  `json:"result,omitempty"`
					StartLocation *struct {
						Alias    string `json:"alias"`
						Yardline int    `json:"yardline"`
					} `json:"start_location,omitempty"`
					EndLocation *struct {
						Alias    string `json:"alias"`
						Yardline int    `json:"yardline"`
					} `json:"end_location,omitempty"`
					Players []struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						Jersey   string `json:"jersey"`
						Position string `json:"position"`
						Role     string `json:"role"`
					} `json:"players,omitempty"`
					Review *struct {
						Result   string `json:"result"`
						Type     string `json:"type"`
						Reversed bool   `json:"reversed,omitempty"`
					} `json:"review,omitempty"`
				} `json:"details,omitempty"`
				EventType string `json:"event_type,omitempty"`
			} `json:"events,omitempty"`
			Inside20      bool `json:"inside_20,omitempty"`
			ScoringDrive  bool `json:"scoring_drive,omitempty"`
			PatSuccessful bool `json:"pat_successful,omitempty"`
		} `json:"pbp,omitempty"`
	} `json:"periods,omitempty"`
	Comment string `json:"_comment"`
}
