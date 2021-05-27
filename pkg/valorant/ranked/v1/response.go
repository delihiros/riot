package v1

type LeaderboardDto struct {
	shard        string       `json:"shard"`
	actId        string       `json:"actId"`
	totalPlayers int          `json:"totalPlayers"`
	players      []*PlayerDto `json:"players"`
}

type PlayerDto struct {
	PUUID           string `json:"puuid"`
	GameName        string `json:"gameName"`
	TagLine         string `json:"tabLine"`
	LeaderboardRank int    `json:"leaderboardRank"`
	RankedRating    int    `json:"rankedRating"`
	NumberOfWins    int    `json:"numberOfWins"`
}
