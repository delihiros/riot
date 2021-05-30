package v1

type MatchDto struct {
	Metadata *MetadataDto `json:"metadata"`
	Info     *InfoDto     `json:"metadata"`
}

type MetadataDto struct {
	DataVersion  string   `json:"data_version"`
	MatchID      string   `json:"match_id"`
	Participants []string `json:"participants"`
}

type InfoDto struct {
	GameMode         string       `json:"game_mode"`
	GameType         string       `json:"game_type"`
	GameStartTimeUTC string       `json:"game_start_time_utc"` // TODO: should be converted to time.Time... or not?
	GameVersion      string       `json:"game_version"`
	Players          []*PlayerDto `json:"players"`
	TotalTurnCount   int          `json:"total_turn_count"`
}

type PlayerDto struct {
	PUUID       string   `json:"puuid"`
	DeckID      string   `json:"deck_id"`
	DeckCode    string   `json:"deck_code"`
	Factions    []string `json:"factions"`
	GameOutcome string   `json:"game_outcome"`
	OrderOfPlay int      `json:"order_of_play"`
}
