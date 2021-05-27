package v1

type MatchDto struct {
	MatchInfo    *MatchInfoDto
	Players      []*PlayerDto
	Coaches      []*CoachDto
	Teams        []*TeamDto
	RoundResults []*RoundResultDto
}

type MatchInfoDto struct {
	MatchID            string `json:"matchId"`
	MapID              string `json:"mapId"`
	GameLengthMillis   int    `json:"gameLengthMillis"`
	GameStartMillis    int    `json:"gameStartMillis"`
	ProvisioningFlowID string `json:"provisioningFlowId"`
	IsCompleted        bool   `json:"isCompleted"`
	CustomGameName     string `json:"customGameName"`
	QueueID            string `json:"queueId"`
	GameMode           string `json:"gameMode"`
	IsRanked           bool   `json:"isRanked"`
	SeasonID           string `json:"seasonId"`
}

type PlayerDto struct {
	PUUID           string          `json:"puuid"`
	GameName        string          `json:"gameName"`
	TagLine         string          `json:"tagLine"`
	TeamID          string          `json:"teamId"`
	PartyID         string          `json:"partyId"`
	CharacterID     string          `json:"characterId"`
	Stats           *PlayerStatsDto `json:"stats"`
	CompetitiveTier int             `json:"competitiveTier"`
	PlayerCard      string          `json:"playerCard"`
	PlayerTitle     string          `json:"playerTitle"`
}

type PlayerStatsDto struct {
	Score          int              `json:"score"`
	RoundsPlayed   int              `json:"roundsPlayed"`
	Kills          int              `json:"kills"`
	Deaths         int              `json:"deaths"`
	Assists        int              `json:"assists"`
	PlaytimeMillis int              `json:"playtimeMillis"`
	AbilityCasts   *AbilityCastsDto `json:"abilityCasts"`
}

type AbilityCastsDto struct {
	GrenadeCasts  int `json:"grenadeCasts"`
	Ability1Casts int `json:"ability1Casts"`
	Ability2Casts int `json:"ability2Casts"`
	UltimateCasts int `json:"ultimateCasts"`
}

type CoachDto struct {
	PUUID  string `json:"puuid"`
	TeamID string `json:"teamId"`
}

type TeamDto struct {
	TeamID       string `json:"teamId"`
	Won          bool   `json:"won"`
	RoundsPlayed int    `json:"roundsPlayed"`
	RoundsWon    int    `json:"roundsWon"`
	NumPoints    int    `json:"numPoints"`
}

type RoundResultDto struct {
	RoundNum              int                    `json:"roundNum"`
	RoundResult           string                 `json:"roundResult"`
	RoundCeremony         string                 `json:"roundCeremony"`
	WinningTeam           string                 `json:"winningTeam"`
	BombPlanter           string                 `json:"bombPlanter"`
	BombDefuser           string                 `json:"bombDefuser"`
	PlantRoundTime        int                    `json:"plantRoundTime"`
	PlantPlayerLocations  []*PlayerLocationsDto  `json:"plantPlayerLocations"`
	PlantLocation         *LocationDto           `json:"plantLocation"`
	PlantSite             string                 `json:"plantSite"`
	DefuseRoundTime       int                    `json:"defuseRoundTime"`
	DefusePlayerLocations []*PlayerLocationsDto  `json:"defusePlayerLocations"`
	DefuseLocation        *LocationDto           `json:"defuseLocation"`
	PlayerStats           []*PlayerRoundStatsDto `json:"playerStats"`
	RoundResultCode       string                 `json:"roundResultCode"`
}

type PlayerLocationsDto struct {
	PUUID       string       `json:"puuid"`
	ViewRadians float64      `json:"viewRadians"`
	Location    *LocationDto `json:"location"`
}

type LocationDto struct {
	X int `json:"x"`
	Y int `json:"Y"`
}

type PlayerRoundStatsDto struct {
	PUUID   string       `json:"puuid"`
	Kills   []*KillDto   `json:"kills"`
	Damage  []*DamageDto `json:"damage"`
	Score   int          `json:"score"`
	Economy *EconomyDto  `json:"economy"`
	Ability *AbilityDto  `json:"ability"`
}

type KillDto struct {
	TimeSinceGameStartMillis  int                   `json:"timeSinceGameStartMillis"`
	TimeSinceRoundStartMillis int                   `json:"timeSinceRoundStartMillis"`
	Killer                    string                `json:"killer"`
	Victim                    string                `json:"victim"`
	VictimLocation            *LocationDto          `json:"victimLocation"`
	Assistants                []string              `json:"assistants"`
	PlayerLocations           []*PlayerLocationsDto `json:"playerLocations"`
	FinishingDamage           *FinishingDamageDto   `json:"finishingDamage"`
}

type FinishingDamageDto struct {
	DamageType          string `json:"damageType"`
	DamageItem          string `json:"damageItem"`
	IsSecondaryFireMode bool   `json:"isSecondaryFireMode"`
}

type DamageDto struct {
	Receiver  string `json:"receiver"`
	Damage    int    `json:"damage"`
	Legshots  int    `json:"legshots"`
	Bodyshots int    `json:"bodyshots"`
	Headshots int    `json:"headshots"`
}

type EconomyDto struct {
	LoadoutValue int    `json:"loadoutValue"`
	Weapon       string `json:"weapon"`
	Armor        string `json:"armor"`
	Remaining    int    `json:"remaining"`
	Spent        int    `json:"spent"`
}

type AbilityDto struct {
	GrenadeEffects  string `json:"grenadeEffects"`
	Ability1Effects string `json:"ability1Effects"`
	Ability2Effects string `json:"ability2Effects"`
	UltimateEffects string `json:"ultimateEffects"`
}

type MatchListDto struct {
	PUUID   string `json:"puuid"`
	History []*MatchListEntryDto
}

type MatchListEntryDto struct {
	MatchID             string `json:"matchID"`
	GameStartTimeMillis int    `json:"gameStartTimeMillis"`
	TeamID              string `json:"teamId"`
}
