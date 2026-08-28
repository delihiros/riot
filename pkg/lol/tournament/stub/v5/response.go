package v5

type TournamentCodeV5DTO struct {
	Code         string   `json:"code"`
	LobbyName    string   `json:"lobbyName"`
	Metadata     string   `json:"metaData"`
	Password     string   `json:"password"`
	TeamSize     int      `json:"teamSize"`
	ProviderID   int      `json:"providerId"`
	PickType     string   `json:"pickType"`
	TournamentID int      `json:"tournamentId"`
	ID           int      `json:"id"`
	Region       string   `json:"region"`
	Map          string   `json:"map"`
	Participants []string `json:"participants"`
}

type LobbyEventV5DTOWrapper struct {
	EventList []LobbyEventV5DTO `json:"eventList"`
}

type LobbyEventV5DTO struct {
	Timestamp string `json:"timestamp"`
	EventType string `json:"eventType"`
	PUUID     string `json:"puuid"`
}
