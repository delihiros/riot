package v1

import pcmatch "github.com/delihiros/riot/pkg/valorant/match/v1"

type PlayerStatsDto = pcmatch.PlayerStatsDto
type CoachDto = pcmatch.CoachDto
type TeamDto = pcmatch.TeamDto
type PlayerLocationsDto = pcmatch.PlayerLocationsDto
type LocationDto = pcmatch.LocationDto
type PlayerRoundStatsDto = pcmatch.PlayerRoundStatsDto
type MatchListDto = pcmatch.MatchListDto
type RecentMatchesDto = pcmatch.RecentMatchesDto

type MatchDto struct {
	MatchInfo    *MatchInfoDto     `json:"matchInfo"`
	Players      []*PlayerDto      `json:"players"`
	Coaches      []*CoachDto       `json:"coaches"`
	Teams        []*TeamDto        `json:"teams"`
	RoundResults []*RoundResultDto `json:"roundResults"`
}

type MatchInfoDto struct {
	MatchID            string `json:"matchId"`
	MapID              string `json:"mapId"`
	GameLengthMillis   int    `json:"gameLengthMillis"`
	GameStartMillis    int64  `json:"gameStartMillis"`
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
