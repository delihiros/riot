package v5

type TournamentCodeParametersV5 struct {
	AllowedParticipants *[]string `json:"allowedParticipants,omitempty"`
	Metadata            *string   `json:"metadata,omitempty"`
	TeamSize            int       `json:"teamSize"`
	PickType            string    `json:"pickType"`
	MapType             string    `json:"mapType"`
	SpectatorType       string    `json:"spectatorType"`
	EnoughPlayers       bool      `json:"enoughPlayers"`
}

type ProviderRegistrationParametersV5 struct {
	Region string `json:"region"`
	URL    string `json:"url"`
}

type TournamentRegistrationParametersV5 struct {
	ProviderID int     `json:"providerId"`
	Name       *string `json:"name,omitempty"`
}

type TournamentCodeUpdateParametersV5 struct {
	AllowedParticipants *[]string `json:"allowedParticipants,omitempty"`
	PickType            *string   `json:"pickType,omitempty"`
	MapType             *string   `json:"mapType,omitempty"`
	SpectatorType       *string   `json:"spectatorType,omitempty"`
}
