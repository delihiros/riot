package v1

type State string

type Tracking string

type ChallengeConfigInfoDto struct {
	ID             int64                        `json:"id"`
	LocalizedNames map[string]map[string]string `json:"localizedNames"`
	State          State                        `json:"state"`
	Tracking       Tracking                     `json:"tracking"`
	StartTimestamp int64                        `json:"startTimestamp"`
	EndTimestamp   int64                        `json:"endTimestamp"`
	Leaderboard    bool                         `json:"leaderboard"`
	Thresholds     map[string]float64           `json:"thresholds"`
}

type ApexPlayerInfoDto struct {
	PUUID    string  `json:"puuid"`
	Value    float64 `json:"value"`
	Position int     `json:"position"`
}

type PlayerInfoDto struct {
	Challenges     []ChallengeInfoDto           `json:"challenges"`
	Preferences    PlayerClientPreferencesDto   `json:"preferences"`
	TotalPoints    ChallengePointDto            `json:"totalPoints"`
	CategoryPoints map[string]ChallengePointDto `json:"categoryPoints"`
}

type ChallengeInfoDto struct {
	Percentile     float64 `json:"percentile"`
	PlayersInLevel int     `json:"playersInLevel"`
	AchievedTime   int64   `json:"achievedTime"`
	Value          float64 `json:"value"`
	ChallengeID    int64   `json:"challengeId"`
	Level          string  `json:"level"`
	Position       int     `json:"position"`
}

type PlayerClientPreferencesDto struct {
	BannerAccent             string   `json:"bannerAccent"`
	Title                    string   `json:"title"`
	ChallengeIDs             []string `json:"challengeIds"`
	CrestBorder              string   `json:"crestBorder"`
	PrestigeCrestBorderLevel int      `json:"prestigeCrestBorderLevel"`
}

type ChallengePointDto struct {
	Level      string `json:"level"`
	Current    int64  `json:"current"`
	Max        int64  `json:"max"`
	Precentile int64  `json:"precentile"`
}
