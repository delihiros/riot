package v1

type PlayerDto struct {
	PUUID    string `json:"puuid"`
	TeamID   string `json:"teamId"`
	Position string `json:"position"`
	Role     string `json:"role"`
}

type TeamDto struct {
	ID           string      `json:"id"`
	TournamentID int         `json:"tournamentId"`
	Name         string      `json:"name"`
	IconID       int         `json:"iconId"`
	Tier         int         `json:"tier"`
	Captain      string      `json:"captain"`
	Abbreviation string      `json:"abbreviation"`
	Players      []PlayerDto `json:"players"`
}

type TournamentDto struct {
	ID               int                  `json:"id"`
	ThemeID          int                  `json:"themeId"`
	NameKey          string               `json:"nameKey"`
	NameKeySecondary string               `json:"nameKeySecondary"`
	Schedule         []TournamentPhaseDto `json:"schedule"`
}

type TournamentPhaseDto struct {
	ID               int   `json:"id"`
	RegistrationTime int64 `json:"registrationTime"`
	StartTime        int64 `json:"startTime"`
	Cancelled        bool  `json:"cancelled"`
}
