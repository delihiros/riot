package v5

type MatchDto struct {
	Metadata MetadataDto `json:"metadata"`
	Info     InfoDto     `json:"info"`
}

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type InfoDto struct {
	EndOfGameResult    string           `json:"endOfGameResult"`
	GameCreation       int64            `json:"gameCreation"`
	GameDuration       int64            `json:"gameDuration"`
	GameEndTimestamp   int64            `json:"gameEndTimestamp"`
	GameID             int64            `json:"gameId"`
	GameMode           string           `json:"gameMode"`
	GameName           string           `json:"gameName"`
	GameStartTimestamp int64            `json:"gameStartTimestamp"`
	GameType           string           `json:"gameType"`
	GameVersion        string           `json:"gameVersion"`
	MapID              int              `json:"mapId"`
	Participants       []ParticipantDto `json:"participants"`
	PlatformID         string           `json:"platformId"`
	QueueID            int              `json:"queueId"`
	Teams              []TeamDto        `json:"teams"`
	TournamentCode     string           `json:"tournamentCode"`
}

type TeamDto struct {
	Bans       []BanDto      `json:"bans"`
	Objectives ObjectivesDto `json:"objectives"`
	TeamID     int           `json:"teamId"`
	Win        bool          `json:"win"`
}

type BanDto struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type ObjectivesDto struct {
	Baron      ObjectiveDto `json:"baron"`
	Champion   ObjectiveDto `json:"champion"`
	Dragon     ObjectiveDto `json:"dragon"`
	Horde      ObjectiveDto `json:"horde"`
	Inhibitor  ObjectiveDto `json:"inhibitor"`
	RiftHerald ObjectiveDto `json:"riftHerald"`
	Tower      ObjectiveDto `json:"tower"`
}

type ObjectiveDto struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

type PerksDto struct {
	StatPerks PerkStatsDto   `json:"statPerks"`
	Styles    []PerkStyleDto `json:"styles"`
}

type PerkStatsDto struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type PerkStyleDto struct {
	Description string                  `json:"description"`
	Selections  []PerkStyleSelectionDto `json:"selections"`
	Style       int                     `json:"style"`
}

type PerkStyleSelectionDto struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type MissionsDto struct {
	PlayerScore0  int `json:"playerScore0"`
	PlayerScore1  int `json:"playerScore1"`
	PlayerScore2  int `json:"playerScore2"`
	PlayerScore3  int `json:"playerScore3"`
	PlayerScore4  int `json:"playerScore4"`
	PlayerScore5  int `json:"playerScore5"`
	PlayerScore6  int `json:"playerScore6"`
	PlayerScore7  int `json:"playerScore7"`
	PlayerScore8  int `json:"playerScore8"`
	PlayerScore9  int `json:"playerScore9"`
	PlayerScore10 int `json:"playerScore10"`
	PlayerScore11 int `json:"playerScore11"`
}

type ReplayDTO struct {
	Total         int      `json:"total"`
	MatchFileURLs []string `json:"matchFileURLs"`
}
