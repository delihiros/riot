package v4

type SummonerDTO struct {
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	PUUID         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}
