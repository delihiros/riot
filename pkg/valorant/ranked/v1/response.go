package v1

type LeaderboardDto struct {
	Shard        string       `json:"shard"`
	ActID        string       `json:"actId"`
	TotalPlayers int64        `json:"totalPlayers"`
	Players      []*PlayerDto `json:"players"`
}

type PlayerDto struct {
	PUUID           string `json:"puuid"`
	GameName        string `json:"gameName"`
	TagLine         string `json:"tagLine"`
	LeaderboardRank int64  `json:"leaderboardRank"`
	RankedRating    int64  `json:"rankedRating"`
	NumberOfWins    int64  `json:"numberOfWins"`
}
