package v5

type TournamentCodeV5DTO struct {
	ID           int      `json:"id"`
	ProviderID   int      `json:"providerId"`
	TournamentID int      `json:"tournamentId"`
	Code         string   `json:"code"`
	Region       string   `json:"region"`
	Map          string   `json:"map"`
	TeamSize     int      `json:"teamSize"`
	Spectators   string   `json:"spectators"`
	PickType     string   `json:"pickType"`
	LobbyName    string   `json:"lobbyName"`
	Password     string   `json:"password"`
	Metadata     string   `json:"metaData"`
	Participants []string `json:"participants"`
}

type TournamentGamesV5 struct {
	StartTime   int64              `json:"startTime"`
	WinningTeam []TournamentTeamV5 `json:"winningTeam"`
	LosingTeam  []TournamentTeamV5 `json:"losingTeam"`
	ShortCode   string             `json:"shortCode"`
	Metadata    string             `json:"metaData"`
	GameID      int64              `json:"gameId"`
	GameName    string             `json:"gameName"`
	GameType    string             `json:"gameType"`
	GameMap     int                `json:"gameMap"`
	GameMode    string             `json:"gameMode"`
	Region      string             `json:"region"`
}

type TournamentTeamV5 struct {
	PUUID string `json:"puuid"`
}

type LobbyEventV5DTOWrapper struct {
	EventList []LobbyEventV5DTO `json:"eventList"`
}

type LobbyEventV5DTO struct {
	Timestamp string `json:"timestamp"`
	EventType string `json:"eventType"`
	PUUID     string `json:"puuid"`
}
