package betgenius

import (
	"encoding/json"
	"time"
)

type Links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

type SeasonsReply struct {
	Total    int `json:"total"`
	Links    `json:"_links"`
	Embedded struct {
		Seasons []*Season `json:"seasons"`
	} `json:"_embedded"`
}

func (s *SeasonsReply) SeasonsToYear() map[int]int {
	idsToYears := make(map[int]int)
	for _, season := range s.Embedded.Seasons {
		idsToYears[season.ID] = season.Year()
	}
	return idsToYears
}

type Season struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	SportID         int    `json:"sportId"`
	SportName       string `json:"sportName"`
	CompetitionID   int    `json:"competitionId"`
	CompetitionName string `json:"competitionName"`
	Locked          string `json:"locked"`
	Seasonproperty  struct {
		StartDate       string `json:"startDate"`
		EndDate         string `json:"endDate"`
		OutrightIsStart bool   `json:"outrightIsStart"`
		FullDescription string `json:"fullDescription"`
		ShortName       string `json:"shortName"`
		LongName        string `json:"longName"`
		SponsorName     string `json:"sponsorName"`
	} `json:"seasonproperty"`
	Updates    int    `json:"updates"`
	Deleted    bool   `json:"deleted"`
	LastUpdate string `json:"lastUpdate"`
	Links      `json:"_links"`
}

func (s *Season) Year() int {
	dateFormat := "2006-01-02 15:04:05"
	startDate, err := time.Parse(dateFormat, s.Seasonproperty.StartDate)
	if err != nil {
		return s.ID
	}

	return startDate.Year()
}

type GamesOfSeason struct {
	Total    int `json:"total"`
	Links    `json:"_links"`
	Embedded struct {
		Fixtures []*Fixture `json:"fixtures"`
	} `json:"_embedded"`
}

type Fixture struct {
	ID                 int    `json:"id"`
	SportID            int    `json:"sportId"`
	Name               string `json:"name"`
	SportName          string `json:"sportName"`
	CompetitionID      string `json:"competitionId"`
	CompetitionName    string `json:"competitionName"`
	SeasonID           int    `json:"seasonId"`
	SeasonName         string `json:"seasonName"`
	RoundID            int    `json:"roundId"`
	RoundName          string `json:"roundName"`
	StartDate          string `json:"startDate"`
	StatusType         string `json:"statusType"`
	Type               string `json:"type"`
	Fixturecompetitors []struct {
		ID                        int `json:"id"`
		Number                    int `json:"number"`
		Fixturecompetitorproperty struct {
		} `json:"fixturecompetitorproperty"`
		Updates    int    `json:"updates"`
		Deleted    bool   `json:"deleted"`
		LastUpdate string `json:"lastUpdate"`
		Links      `json:"_links"`
		Competitor struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			SportID    int    `json:"sportId"`
			Gender     string `json:"gender"`
			Updates    int    `json:"updates"`
			Deleted    bool   `json:"deleted"`
			LastUpdate string `json:"lastUpdate"`
			Links      struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			Type         string `json:"type"`
			Teamproperty struct {
				IsProtected     string `json:"isProtected"`
				ShortName       string `json:"shortName"`
				FullDescription string `json:"fullDescription"`
			} `json:"teamproperty"`
			Playerproperty json.RawMessage `json:"playerproperty"`
		} `json:"competitor"`
	} `json:"fixturecompetitors"`
	Venue struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		RegionID      int    `json:"regionId"`
		VenueType     string `json:"venueType"`
		Venueproperty struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
			Capacity  string `json:"capacity"`
			City      string `json:"city"`
		} `json:"venueproperty"`
		Updates    int    `json:"updates"`
		Deleted    bool   `json:"deleted"`
		LastUpdate string `json:"lastUpdate"`
		Links      `json:"_links"`
	} `json:"venue"`
	Fixtureproperty struct {
		HasLineups   bool   `json:"hasLineups"`
		Timezone     string `json:"timezone"`
		NeutralField bool   `json:"neutralField"`
	} `json:"fixtureproperty"`
	Updates    int    `json:"updates"`
	Deleted    bool   `json:"deleted"`
	LastUpdate string `json:"lastUpdate"`
	Links      `json:"_links"`
}

type Drive struct {
	TeamInPossession string `json:"teamInPossession"`
	IsKickOff        bool   `json:"isKickOff"`
	Plays            []struct {
		Sequence          int  `json:"sequence"`
		DownNumber        *int `json:"downNumber"`
		YardsToGo         *int `json:"yardsToGo"`
		ScrimmageLocation struct {
			ScrimmageYard int    `json:"scrimmageYard"`
			SideOfPitch   string `json:"sideOfPitch"`
		} `json:"scrimmageLocation"`
		Snap        *Snap  `json:"snap"`
		ID          string `json:"id"`
		IsVoid      bool   `json:"isVoid"`
		IsConfirmed bool   `json:"isConfirmed"`
		IsFinished  bool   `json:"isFinished"`
		Period      struct {
			Number int    `json:"number"`
			Type   string `json:"type"`
		} `json:"period"`
		Actions []struct {
			ID       string  `json:"id"`
			Sequence int     `json:"sequence"`
			Team     string  `json:"team"`
			Type     string  `json:"type"`
			SubType  *string `json:"subType"`
			Players  []struct {
				PlayerType string `json:"playerType"`
				ID         string `json:"id"`
			} `json:"players"`
			Yards       int       `json:"yards"`
			YardLine    *YardLine `json:"yardLine"`
			IsNullified bool      `json:"isNullified"`
			ActiveUnits struct {
				Home string `json:"home"`
				Away string `json:"away"`
			} `json:"activeUnits"`
		} `json:"actions"`
		Penalties         []Penalty  `json:"penalties"`
		StartedAtGameTime string     `json:"startedAtGameTime"`
		EndedAtGameTime   string     `json:"endedAtGameTime"`
		StartedAtUtc      *time.Time `json:"startedAtUtc"`
		EndedAtUtc        *time.Time `json:"endedAtUtc"`
		Description       string     `json:"description"`
		SourcePlayID      string     `json:"sourcePlayId"`
	} `json:"plays"`
	ConversionPlays []interface{} `json:"conversionPlays"`
	Score           []*Score      `json:"score"`
	IsFinished      bool          `json:"isFinished"`
}
type GamePbp struct {
	FirstHalf struct {
		Drives   []Drive `json:"drives"`
		Timeouts []struct {
			Team     string `json:"team"`
			GameTime string `json:"gameTime"`
			Period   struct {
				Number int    `json:"number"`
				Type   string `json:"type"`
			} `json:"period"`
			UtcTimestamp *time.Time `json:"utcTimestamp"`
			IsConfirmed  bool       `json:"isConfirmed"`
		} `json:"timeouts"`
		TimeoutsRemaining struct {
			Away int `json:"away"`
			Home int `json:"home"`
		} `json:"timeoutsRemaining"`
		CoinToss struct {
			WinnerTeam  string `json:"winnerTeam"`
			WasDeferred bool   `json:"wasDeferred"`
			AwayChoice  string `json:"awayChoice"`
			HomeChoice  string `json:"homeChoice"`
		} `json:"coinToss"`
	} `json:"firstHalf"`
	SecondHalf struct {
		Drives   []Drive `json:"drives"`
		Timeouts []struct {
			Team     string `json:"team"`
			GameTime string `json:"gameTime"`
			Period   struct {
				Number int    `json:"number"`
				Type   string `json:"type"`
			} `json:"period"`
			UtcTimestamp *time.Time `json:"utcTimestamp"`
			IsConfirmed  bool       `json:"isConfirmed"`
		} `json:"timeouts"`
		TimeoutsRemaining struct {
			Away int `json:"away"`
			Home int `json:"home"`
		} `json:"timeoutsRemaining"`
		CoinToss struct {
			WinnerTeam  string `json:"winnerTeam"`
			WasDeferred bool   `json:"wasDeferred"`
			AwayChoice  string `json:"awayChoice"`
			HomeChoice  string `json:"homeChoice"`
		} `json:"coinToss"`
	} `json:"secondHalf"`
	OvertimePeriods []OvertimePeriod `json:"overtimePeriods"`
	Challenges      []interface{}    `json:"challenges"`
	GameTime        struct {
		Clock          string     `json:"clock"`
		LastUpdatedUtc *time.Time `json:"lastUpdatedUtc"`
		IsRunning      bool       `json:"isRunning"`
	} `json:"gameTime"`
	PlayClock struct {
		Clock          string     `json:"clock"`
		LastUpdatedUtc *time.Time `json:"lastUpdatedUtc"`
		IsRunning      bool       `json:"isRunning"`
	} `json:"playClock"`
	MatchStatus string `json:"matchStatus"`
	Period      struct {
		Number int    `json:"number"`
		Type   string `json:"type"`
	} `json:"period"`
	PeriodWithStatus struct {
		Status      string `json:"status"`
		IsConfirmed bool   `json:"isConfirmed"`
		Number      int    `json:"number"`
		Type        string `json:"type"`
	} `json:"periodWithStatus"`
	Score struct {
		Home        int  `json:"home"`
		Away        int  `json:"away"`
		IsConfirmed bool `json:"isConfirmed"`
	} `json:"score"`
	HomeTeam struct {
		Offensive []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"offensive"`
		Defensive []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"defensive"`
		Special []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"special"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Offensive []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"offensive"`
		Defensive []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"defensive"`
		Special []struct {
			ID       string `json:"id"`
			Position string `json:"position"`
			Side     string `json:"side"`
			Status   string `json:"status"`
		} `json:"special"`
	} `json:"awayTeam"`
	Injuries []struct {
		ID           string     `json:"id"`
		PlayID       *string    `json:"playId"`
		Team         string     `json:"team"`
		PlayerID     string     `json:"playerId"`
		Status       string     `json:"status"`
		IsConfirmed  bool       `json:"isConfirmed"`
		UtcTimestamp *time.Time `json:"utcTimestamp"`
	} `json:"injuries"`
	Comments []struct {
		Value        string     `json:"value"`
		UtcTimestamp *time.Time `json:"utcTimestamp"`
	} `json:"comments"`
	CurrentPossession struct {
		Team           string     `json:"team"`
		LastUpdatedUtc *time.Time `json:"lastUpdatedUtc"`
	} `json:"currentPossession"`
	YardsToEndzone struct {
		Yards          int        `json:"yards"`
		Team           string     `json:"team"`
		LastUpdatedUtc *time.Time `json:"lastUpdatedUtc"`
		IsConfirmed    bool       `json:"isConfirmed"`
	} `json:"yardsToEndzone"`
	Risks struct {
		Touchdown        string `json:"touchdown"`
		OnsideKick       string `json:"onsideKick"`
		FieldGoal        string `json:"fieldGoal"`
		FourthDown       string `json:"fourthDown"`
		Safety           string `json:"safety"`
		Challenge        string `json:"challenge"`
		Penalty          string `json:"penalty"`
		VideoReview      string `json:"videoReview"`
		Turnover         string `json:"turnover"`
		Other            string `json:"other"`
		PlayAboutToStart string `json:"playAboutToStart"`
		Injury           string `json:"injury"`
		BigPlay          string `json:"bigPlay"`
		StatDelay        string `json:"statDelay"`
	} `json:"risks"`
	IsPlayUnderReview   *bool           `json:"isPlayUnderReview"`
	NextPlay            json.RawMessage `json:"nextPlay"`
	Source              string          `json:"source"`
	FixtureID           string          `json:"fixtureId"`
	Sequence            int             `json:"sequence"`
	MessageTimestampUtc *time.Time      `json:"messageTimestampUtc"`
	IsReliable          bool            `json:"isReliable"`
	IsCoverageCancelled bool            `json:"isCoverageCancelled"`
	ReliabilityReasons  struct {
		Heartbeat       string `json:"heartbeat"`
		FeedReliability string `json:"feedReliability"`
		Coverage        string `json:"coverage"`
	} `json:"reliabilityReasons"`
}

type OvertimePeriod struct {
	Period   int     `json:"period"`
	Drives   []Drive `json:"drives"`
	Timeouts []struct {
		Team     string `json:"team"`
		GameTime string `json:"gameTime"`
		Period   struct {
			Number int    `json:"number"`
			Type   string `json:"type"`
		} `json:"period"`
		UtcTimestamp time.Time `json:"utcTimestamp"`
		IsConfirmed  bool      `json:"isConfirmed"`
	} `json:"timeouts"`
	TimeoutsRemaining struct {
		Away int `json:"away"`
		Home int `json:"home"`
	} `json:"timeoutsRemaining"`
	CoinToss struct {
		WinnerTeam  string `json:"winnerTeam"`
		WasDeferred bool   `json:"wasDeferred"`
		AwayChoice  string `json:"awayChoice"`
		HomeChoice  string `json:"homeChoice"`
	} `json:"coinToss"`
}

type Snap struct {
	ID           string     `json:"id"`
	IsConfirmed  bool       `json:"isConfirmed"`
	TimestampUtc *time.Time `json:"timestampUtc"`
}

type YardLine struct {
	Yards       int    `json:"yards"`
	SideOfPitch string `json:"sideOfPitch"`
}

type Score struct {
	Period struct {
		Number int    `json:"number"`
		Type   string `json:"type"`
	} `json:"period"`
	Type         string     `json:"type"`
	Team         string     `json:"team"`
	Points       int        `json:"points"`
	IsConfirmed  bool       `json:"isConfirmed"`
	UtcTimestamp *time.Time `json:"utcTimestamp"`
}

type ConversionPlay struct {
	TeamInPossession string `json:"teamInPossession"`
	Type             string `json:"type"`
	ID               string `json:"id"`
	IsVoid           bool   `json:"isVoid"`
	IsConfirmed      bool   `json:"isConfirmed"`
	IsFinished       bool   `json:"isFinished"`
	Period           struct {
		Number int    `json:"number"`
		Type   string `json:"type"`
	} `json:"period"`
	Actions []struct {
		ID       string  `json:"id"`
		Sequence int     `json:"sequence"`
		Team     string  `json:"team"`
		Type     string  `json:"type"`
		SubType  *string `json:"subType"`
		Players  []struct {
			PlayerType string `json:"playerType"`
			ID         string `json:"id"`
		} `json:"players"`
		Yards       interface{} `json:"yards"`
		YardLine    interface{} `json:"yardLine"`
		IsNullified bool        `json:"isNullified"`
		ActiveUnits struct {
			Home string `json:"home"`
			Away string `json:"away"`
		} `json:"activeUnits"`
	} `json:"actions"`
	Penalties         []Penalty  `json:"penalties"`
	StartedAtGameTime string     `json:"startedAtGameTime"`
	EndedAtGameTime   *time.Time `json:"endedAtGameTime"`
	StartedAtUtc      *time.Time `json:"startedAtUtc"`
	EndedAtUtc        *time.Time `json:"endedAtUtc"`
	Description       string     `json:"description"`
	SourcePlayID      string     `json:"sourcePlayId"`
}

type Penalty struct {
	ID        string `json:"id"`
	Team      string `json:"team"`
	PlayerID  string `json:"playerId"`
	Type      string `json:"type"`
	Outcome   string `json:"outcome"`
	Yards     int    `json:"yards"`
	YardLines []struct {
		Type        string `json:"type"`
		Yards       int    `json:"yards"`
		SideOfPitch string `json:"sideOfPitch"`
	} `json:"yardLines"`
	EnforcementSpot string `json:"enforcementSpot"`
	NextDown        string `json:"nextDown"`
	ActiveUnits     struct {
		Home string `json:"home"`
		Away string `json:"away"`
	} `json:"activeUnits"`
	UtcTimestamp *time.Time `json:"utcTimestamp"`
}
