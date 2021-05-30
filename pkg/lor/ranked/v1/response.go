package v1

type LeaderboardDto struct {
	Players []*PlayerDto `json:"players"`
}

type PlayerDto struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
	LP   int    `json:"lp"`
}
