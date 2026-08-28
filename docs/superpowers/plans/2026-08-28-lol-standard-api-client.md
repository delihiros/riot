# LoL standard platform API client implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add typed clients for the nine standard platform-routed LoL API groups, covering 25 API-key operations.

**Architecture:** Each Riot API group becomes a versioned package below `pkg/lol` and wraps the existing shared `pkg/client.Client`. Methods escape path parameters, encode optional queries, send the configured Riot API key, and decode complete documented response DTOs.

**Tech Stack:** Go 1.26.7, standard library `encoding/json`, `net/http`, `net/url`, and the repository's existing test transport pattern.

---

**Spec:** `docs/superpowers/specs/2026-08-28-lol-api-client-design.md`

**Official contracts:** Riot Developer Portal API details for `champion-mastery-v4`, `champion-v3`, `clash-v1`, `league-exp-v4`, `league-v4`, `lol-challenges-v1`, `lol-status-v4`, `spectator-v5`, and `summoner-v4`, checked 2026-08-28.

**Required workflow:** Use @superpowers:test-driven-development for every production method. Every test must fail because the method or DTO is absent before implementation begins. In Step 1 of every package task, also add one malformed-JSON case that expects a decode error and one non-2xx case that expects the shared `HTTPError`; these cases must call a new package method. Before every GREEN run and task commit, run `gofmt -w` on that task's Go files.

### Task 1: Champion Mastery v4

**Files:**
- Create: `pkg/lol/championmastery/v4/client.go`
- Create: `pkg/lol/championmastery/v4/response.go`
- Create: `pkg/lol/championmastery/v4/client_test.go`

- [ ] **Step 1: Write failing endpoint tests**

Add tests for these exact public methods:

```go
GetAllChampionMasteriesByPUUID(region, puuid string) ([]ChampionMasteryDto, error)
GetChampionMasteryByPUUID(region, puuid string, championID int64) (*ChampionMasteryDto, error)
GetTopChampionMasteriesByPUUID(region, puuid string, count *int) ([]ChampionMasteryDto, error)
GetChampionMasteryScoreByPUUID(region, puuid string) (int, error)
```

Assert all four official paths, escaped PUUIDs, optional `count`, `X-Riot-Token`, list/object/primitive decoding, and the absence of a query when `count` is nil. Use a complete fixture for every documented field.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/championmastery/v4`

Expected: build failure because `Client`, methods, and DTOs do not exist.

- [ ] **Step 3: Implement the minimal client and DTOs**

Create `Client`, `New`, the four methods, and these exact DTO shapes:

```text
ChampionMasteryDto: puuid string; championPointsUntilNextLevel int64; chestGranted bool; championId int64; lastPlayTime int64; championLevel int; championPoints int; championPointsSinceLastLevel int64; markRequiredForNextLevel int; championSeasonMilestone int; nextSeasonMilestone NextSeasonMilestonesDto; tokensEarned int; milestoneGrades []string
NextSeasonMilestonesDto: requireGradeCounts map[string]int; rewardMarks int; bonus bool; rewardConfig RewardConfigDto
RewardConfigDto: rewardValue string; rewardType string; maximumReward int
```

Use exported Go field names and the exact lower-camel JSON tags shown above.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/championmastery/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/championmastery/v4 && git commit -m "feat(lol): add champion mastery v4 client"`

### Task 2: Champion v3

**Files:**
- Create: `pkg/lol/champion/v3/client.go`
- Create: `pkg/lol/champion/v3/response.go`
- Create: `pkg/lol/champion/v3/client_test.go`

- [ ] **Step 1: Write a failing rotation test**

Test `GetChampionInfo(region string) (*ChampionInfo, error)` against `GET /lol/platform/v3/champion-rotations`. Assert the API-key header and decode both complete fields:

```go
type ChampionInfo struct {
    NewPlayer []int `json:"newplayer"`
    SR        []int `json:"sr"`
}
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/champion/v3`

Expected: build failure for the missing client.

- [ ] **Step 3: Implement the endpoint and response**

Use `SimpleGet`, decode `ChampionInfo`, and return decoding errors unchanged.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/champion/v3`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/champion/v3 && git commit -m "feat(lol): add champion v3 client"`

### Task 3: Clash v1

**Files:**
- Create: `pkg/lol/clash/v1/client.go`
- Create: `pkg/lol/clash/v1/response.go`
- Create: `pkg/lol/clash/v1/client_test.go`

- [ ] **Step 1: Write failing tests for all five operations**

Test these signatures and paths:

```go
GetPlayersByPUUID(region, puuid string) ([]PlayerDto, error)
GetTeamByID(region, teamID string) (*TeamDto, error)
GetTournaments(region string) ([]TournamentDto, error)
GetTournamentByTeam(region, teamID string) (*TournamentDto, error)
GetTournamentByID(region string, tournamentID int) (*TournamentDto, error)
```

Use full response fixtures for:

```text
PlayerDto: puuid string; position string; role string
TeamDto: id string; tournamentId int; name string; iconId int; tier int; captain string; abbreviation string; players []PlayerDto
TournamentDto: id int; themeId int; nameKey string; nameKeySecondary string; schedule []TournamentPhaseDto
TournamentPhaseDto: id int; registrationTime int64; startTime int64; cancelled bool
```

Assert that team IDs and PUUIDs containing `/` are path-escaped.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/clash/v1`

Expected: build failure for missing methods.

- [ ] **Step 3: Implement client and DTOs**

Implement only the five documented GET methods and exact JSON tags.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/clash/v1`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/clash/v1 && git commit -m "feat(lol): add clash v1 client"`

### Task 4: League Exp v4

**Files:**
- Create: `pkg/lol/leagueexp/v4/client.go`
- Create: `pkg/lol/leagueexp/v4/response.go`
- Create: `pkg/lol/leagueexp/v4/client_test.go`

- [ ] **Step 1: Write a failing entries test**

Test:

```go
GetLeagueEntries(region, queue, tier, division string, page *int) ([]LeagueEntryDTO, error)
```

Assert `/lol/league-exp/v4/entries/{queue}/{tier}/{division}`, escaping for all segments, omitted/present `page`, and complete set-as-slice decoding.

Response fields:

```text
LeagueEntryDTO: leagueId string; summonerId string; puuid string; queueType string; tier string; rank string; leaguePoints int; wins int; losses int; hotStreak bool; veteran bool; freshBlood bool; inactive bool; miniSeries MiniSeriesDTO
MiniSeriesDTO: losses int; progress string; target int; wins int
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/leagueexp/v4`

Expected: build failure.

- [ ] **Step 3: Implement endpoint and DTOs**

Encode `page` only when non-nil and decode Riot's `Set[LeagueEntryDTO]` JSON array into `[]LeagueEntryDTO`.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/leagueexp/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/leagueexp/v4 && git commit -m "feat(lol): add league exp v4 client"`

### Task 5: League v4

**Files:**
- Create: `pkg/lol/league/v4/client.go`
- Create: `pkg/lol/league/v4/response.go`
- Create: `pkg/lol/league/v4/client_test.go`

- [ ] **Step 1: Write failing tests for all five operations**

Test:

```go
GetChallengerLeague(region, queue string) (*LeagueListDTO, error)
GetLeagueEntriesByPUUID(region, puuid string) ([]LeagueEntryDTO, error)
GetLeagueEntries(region, queue, tier, division string, page *int) ([]LeagueEntryDTO, error)
GetGrandmasterLeague(region, queue string) (*LeagueListDTO, error)
GetMasterLeague(region, queue string) (*LeagueListDTO, error)
```

Assert official paths, optional `page`, segment escaping, API-key header, and complete decoding.

Response fields:

```text
LeagueListDTO: leagueId string; entries []LeagueItemDTO; tier string; name string; queue string
LeagueItemDTO: freshBlood bool; wins int; miniSeries MiniSeriesDTO; inactive bool; veteran bool; hotStreak bool; rank string; leaguePoints int; losses int; puuid string
LeagueEntryDTO: leagueId string; puuid string; queueType string; tier string; rank string; leaguePoints int; wins int; losses int; hotStreak bool; veteran bool; freshBlood bool; inactive bool; miniSeries MiniSeriesDTO
MiniSeriesDTO: losses int; progress string; target int; wins int
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/league/v4`

Expected: build failure.

- [ ] **Step 3: Implement client and DTOs**

Use `[]LeagueEntryDTO` for both documented set responses.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/league/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/league/v4 && git commit -m "feat(lol): add league v4 client"`

### Task 6: Challenges v1

**Files:**
- Create: `pkg/lol/challenges/v1/client.go`
- Create: `pkg/lol/challenges/v1/response.go`
- Create: `pkg/lol/challenges/v1/client_test.go`

- [ ] **Step 1: Write failing tests for all six operations**

Test:

```go
GetAllChallengeConfigs(region string) ([]ChallengeConfigInfoDto, error)
GetAllChallengePercentiles(region string) (map[int64]map[int]map[string]float64, error)
GetChallengeConfigs(region string, challengeID int64) (*ChallengeConfigInfoDto, error)
GetChallengeLeaderboards(region string, challengeID int64, level string, limit *int) ([]ApexPlayerInfoDto, error)
GetChallengePercentiles(region string, challengeID int64) (map[string]float64, error)
GetPlayerData(region, puuid string) (*PlayerInfoDto, error)
```

Assert numeric JSON object keys decode correctly, `limit` is omitted when nil, and PUUID/level are escaped.

Response fields:

```text
ChallengeConfigInfoDto: id int64; localizedNames map[string]map[string]string; state State; tracking Tracking; startTimestamp int64; endTimestamp int64; leaderboard bool; thresholds map[string]float64
State: string alias
Tracking: string alias
ApexPlayerInfoDto: puuid string; value float64; position int
PlayerInfoDto: challenges []ChallengeInfoDto; preferences PlayerClientPreferencesDto; totalPoints ChallengePointDto; categoryPoints map[string]ChallengePointDto
ChallengeInfoDto: percentile float64; playersInLevel int; achievedTime int64; value float64; challengeId int64; level string; position int
PlayerClientPreferencesDto: bannerAccent string; title string; challengeIds []string; crestBorder string; prestigeCrestBorderLevel int
ChallengePointDto: level string; current int64; max int64; precentile int64
```

Preserve Riot's documented `precentile` spelling in the JSON tag and Go field name so no undocumented alias is invented.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/challenges/v1`

Expected: build failure.

- [ ] **Step 3: Implement client and DTOs**

Use `url.Values` for `limit` and standard JSON decoding for numeric map keys.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/challenges/v1`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/challenges/v1 && git commit -m "feat(lol): add challenges v1 client"`

### Task 7: Status v4

**Files:**
- Create: `pkg/lol/status/v4/client.go`
- Create: `pkg/lol/status/v4/response.go`
- Create: `pkg/lol/status/v4/client_test.go`

- [ ] **Step 1: Write a failing platform-data test**

Test `GetPlatformData(region string) (*PlatformDataDto, error)` and the exact path `/lol/status/v4/platform-data`.

Use a fixture covering:

```text
PlatformDataDto: id string; name string; locales []string; maintenances []StatusDto; incidents []StatusDto
StatusDto: id int; maintenance_status string; incident_severity string; titles []ContentDto; updates []UpdateDto; created_at string; archive_at string; updated_at string; platforms []string
ContentDto: locale string; content string
UpdateDto: id int; author string; publish bool; publish_locations []string; translations []ContentDto; created_at string; updated_at string
```

Use idiomatic exported Go field names and preserve all snake-case JSON tags.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/status/v4`

Expected: build failure.

- [ ] **Step 3: Implement endpoint and DTOs**

Use `SimpleGet` and exact response tags.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/status/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/status/v4 && git commit -m "feat(lol): add status v4 client"`

### Task 8: Spectator v5

**Files:**
- Create: `pkg/lol/spectator/v5/client.go`
- Create: `pkg/lol/spectator/v5/response.go`
- Create: `pkg/lol/spectator/v5/client_test.go`

- [ ] **Step 1: Write a failing active-game test**

Test `GetCurrentGameInfoByPUUID(region, puuid string) (*CurrentGameInfo, error)` and `/lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}`.

Response fields:

```text
CurrentGameInfo: gameId int64; gameType string; gameStartTime int64; mapId int64; gameLength int64; platformId string; gameMode string; bannedChampions []BannedChampion; gameQueueConfigId int64; observers Observer; participants []CurrentGameParticipant
BannedChampion: pickTurn int; championId int64; teamId int64
Observer: encryptionKey string
CurrentGameParticipant: championId int64; perks Perks; profileIconId int64; bot bool; teamId int64; puuid string; spell1Id int64; spell2Id int64; gameCustomizationObjects []GameCustomizationObject
Perks: perkIds []int64; perkStyle int64; perkSubStyle int64
GameCustomizationObject: category string; content string
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/spectator/v5`

Expected: build failure.

- [ ] **Step 3: Implement endpoint and DTOs**

Escape the PUUID and decode the complete fixture.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/spectator/v5`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/spectator/v5 && git commit -m "feat(lol): add spectator v5 client"`

### Task 9: Summoner v4 API-key endpoint

**Files:**
- Create: `pkg/lol/summoner/v4/client.go`
- Create: `pkg/lol/summoner/v4/response.go`
- Create: `pkg/lol/summoner/v4/client_test.go`

- [ ] **Step 1: Write a failing PUUID lookup test**

Test:

```go
GetByPUUID(region, puuid string) (*SummonerDTO, error)
```

Assert `/lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}`, escaping, API-key header, and all fields:

```text
SummonerDTO: profileIconId int; revisionDate int64; puuid string; summonerLevel int64
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/summoner/v4`

Expected: build failure.

- [ ] **Step 3: Implement endpoint and DTO**

Leave the RSO `/me` method for the Match/RSO plan so this task remains independently shippable.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/summoner/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/summoner/v4 && git commit -m "feat(lol): add summoner v4 client"`

### Task 10: Standard API integration verification

**Files:**
- Verify only: all files created in Tasks 1-9

- [ ] **Step 1: Count implemented operations**

Confirm the nine packages expose exactly 25 public endpoint methods: 4 + 1 + 5 + 1 + 5 + 6 + 1 + 1 + 1. The second Summoner operation is explicitly supplied by the Match/RSO plan.

- [ ] **Step 2: Run the full test suite**

Run: `GOTOOLCHAIN=go1.26.7 go test ./...`

Expected: PASS.

- [ ] **Step 3: Run race and static analysis**

Run: `GOTOOLCHAIN=go1.26.7 go test -race ./...`

Run: `GOTOOLCHAIN=go1.26.7 go vet ./...`

Expected: both PASS.

- [ ] **Step 4: Check formatting and worktree**

Run: `gofmt -l pkg/lol`

Expected: no output; every task formatted its files before committing.

Run: `git diff --check`

Expected: no formatting or whitespace errors.
