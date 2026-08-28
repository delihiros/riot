# LoL Tournament API client implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Tournament Stub v5 and Tournament v5 clients, the shared authenticated write helpers they require, and final verification for all 45 LoL operations.

**Architecture:** `pkg/client` remains the only code that can read the configured Riot API key and gains authenticated JSON POST/PUT helpers. Tournament packages provide typed request and response DTOs and use those helpers for write operations while retaining the existing shared error and security behavior.

**Tech Stack:** Go 1.26.7, standard library `bytes`, `encoding/json`, `net/http`, `net/url`, and the existing shared client.

---

**Spec:** `docs/superpowers/specs/2026-08-28-lol-api-client-design.md`

**Depends on:** The standard and Match/RSO plans are complete.

**Official contracts:** Riot Developer Portal API details for `tournament-stub-v5` and `tournament-v5`, checked 2026-08-28.

**Required workflow:** Use @superpowers:test-driven-development for every shared helper and endpoint. Tournament request fixtures must distinguish omitted values from explicit empty arrays and zero values. In Step 1 of every Tournament package task, also add one malformed-JSON case that expects a decode error and one non-2xx case that expects the shared `HTTPError`; these cases must call a new package method. Before every GREEN run and task commit, run `gofmt -w` on that task's Go files.

### Task 1: Authenticated shared POST and PUT helpers

**Files:**
- Modify: `pkg/client/client.go`
- Modify: `pkg/client/client_test.go`

- [ ] **Step 1: Write failing authenticated write tests**

Add tests for the desired API:

```go
SimplePost(region, path string, body []byte) ([]byte, error)
SimplePut(region, path string, body []byte) ([]byte, error)
PutWithRegionAndHeaders(region, path string, params map[string]string, body []byte) ([]byte, error)
```

Assert method, Riot host, query preservation, exact body, `X-Riot-Token`, `Content-Type: application/json`, redirect rejection, and invalid routing/path rejection before any request is sent. For `SimplePut` with a nil body, assert an empty request body rather than JSON `null`.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/client`

Expected: build failure because the methods do not exist.

- [ ] **Step 3: Refactor the existing POST path and implement PUT**

Keep the public surface limited to the three methods above. Refactor `PostWithRegionAndHeaders` and the new PUT method through one private helper:

```go
func (c *Client) requestWithRegionAndHeaders(method, region, path string, params map[string]string, body []byte) ([]byte, error)
```

Have `SimplePost` and `SimplePut` construct these headers internally:

```go
map[string]string{
    "X-Riot-Token": c.apiKey,
    "Content-Type": "application/json",
}
```

Do not expose the API key or a generic public arbitrary-method function.

- [ ] **Step 4: Run tests and verify GREEN**

Run: `gofmt -w pkg/client && GOTOOLCHAIN=go1.26.7 go test ./pkg/client`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/client && git commit -m "feat(client): add authenticated JSON write helpers"`

### Task 2: Tournament Stub v5

**Files:**
- Create: `pkg/lol/tournament/stub/v5/client.go`
- Create: `pkg/lol/tournament/stub/v5/request.go`
- Create: `pkg/lol/tournament/stub/v5/response.go`
- Create: `pkg/lol/tournament/stub/v5/client_test.go`

- [ ] **Step 1: Write failing tests for all five operations**

Test:

```go
CreateTournamentCode(region string, tournamentID int64, count *int, params TournamentCodeParametersV5) ([]string, error)
GetTournamentCode(region, tournamentCode string) (*TournamentCodeV5DTO, error)
GetLobbyEventsByCode(region, tournamentCode string) (*LobbyEventV5DTOWrapper, error)
RegisterProviderData(region string, params ProviderRegistrationParametersV5) (int, error)
RegisterTournament(region string, params TournamentRegistrationParametersV5) (int, error)
```

Assert POST/GET methods, exact `/lol/tournament-stub/v5` paths, escaped tournament codes, `tournamentId`, optional `count`, configured API-key header, exact JSON bodies, and list/object/integer decoding.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/tournament/stub/v5`

Expected: build failure.

- [ ] **Step 3: Implement request DTOs**

The POST bodies are required, but Riot marks `allowedParticipants`, `metadata`, and tournament `name` as optional. Use pointers with `omitempty` for those fields and value fields for the required fields:

```go
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
```

Test nil optional pointers are absent, a pointer to an empty participant slice encodes `"allowedParticipants":[]`, and a pointer to an empty string encodes the corresponding empty string.

- [ ] **Step 4: Implement response DTOs and client methods**

Create:

```text
TournamentCodeV5DTO: code string; lobbyName string; metaData string; password string; teamSize int; providerId int; pickType string; tournamentId int; id int; region string; map string; participants []string
LobbyEventV5DTOWrapper: eventList []LobbyEventV5DTO
LobbyEventV5DTO: timestamp string; eventType string; puuid string
```

Preserve `metaData` exactly in the response JSON tag and `metadata` exactly in the request JSON tag.

- [ ] **Step 5: Run package tests and verify GREEN**

Run: `gofmt -w pkg/lol/tournament/stub/v5 && GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/tournament/stub/v5`

Expected: PASS.

- [ ] **Step 6: Commit**

Run: `git add pkg/lol/tournament/stub/v5 && git commit -m "feat(lol): add tournament stub v5 client"`

### Task 3: Tournament v5

**Files:**
- Create: `pkg/lol/tournament/v5/client.go`
- Create: `pkg/lol/tournament/v5/request.go`
- Create: `pkg/lol/tournament/v5/response.go`
- Create: `pkg/lol/tournament/v5/client_test.go`

- [ ] **Step 1: Write failing tests for all seven operations**

Test:

```go
CreateTournamentCode(region string, tournamentID int64, count *int, params TournamentCodeParametersV5) ([]string, error)
GetTournamentCode(region, tournamentCode string) (*TournamentCodeV5DTO, error)
UpdateCode(region, tournamentCode string, params *TournamentCodeUpdateParametersV5) error
GetGames(region, tournamentCode string) ([]TournamentGamesV5, error)
GetLobbyEventsByCode(region, tournamentCode string) (*LobbyEventV5DTOWrapper, error)
RegisterProviderData(region string, params ProviderRegistrationParametersV5) (int, error)
RegisterTournament(region string, params TournamentRegistrationParametersV5) (int, error)
```

Assert the exact `/lol/tournament/v5` methods and paths. For `UpdateCode`, test both nil parameters producing an empty body and selected parameters producing only selected JSON properties. Test a pointer to an empty participant slice encodes `"allowedParticipants":[]`.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/tournament/v5`

Expected: build failure.

- [ ] **Step 3: Implement request DTOs**

Use the same three POST request DTO shapes and optional-field rules as Tournament Stub. Add the partial update DTO:

```go
type TournamentCodeUpdateParametersV5 struct {
    AllowedParticipants *[]string `json:"allowedParticipants,omitempty"`
    PickType            *string   `json:"pickType,omitempty"`
    MapType             *string   `json:"mapType,omitempty"`
    SpectatorType       *string   `json:"spectatorType,omitempty"`
}
```

When `params == nil`, call `SimplePut` with a nil body. Otherwise JSON-marshal the DTO.

- [ ] **Step 4: Implement response DTOs**

Create:

```text
TournamentCodeV5DTO: id int; providerId int; tournamentId int; code string; region string; map string; teamSize int; spectators string; pickType string; lobbyName string; password string; metaData string; participants []string
TournamentGamesV5: startTime int64; winningTeam []TournamentTeamV5; losingTeam []TournamentTeamV5; shortCode string; metaData string; gameId int64; gameName string; gameType string; gameMap int; gameMode string; region string
TournamentTeamV5: puuid string
LobbyEventV5DTOWrapper: eventList []LobbyEventV5DTO
LobbyEventV5DTO: timestamp string; eventType string; puuid string
```

Use `[]TournamentGamesV5` for Riot's documented set response.

- [ ] **Step 5: Implement all seven methods and verify GREEN**

Run: `gofmt -w pkg/lol/tournament/v5 && GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/tournament/v5`

Expected: PASS.

- [ ] **Step 6: Commit**

Run: `git add pkg/lol/tournament/v5 && git commit -m "feat(lol): add tournament v5 client"`

### Task 4: Documentation and complete LoL verification

**Files:**
- Modify: `README.md`
- Verify: every file below `pkg/lol`

- [ ] **Step 1: Verify the endpoint inventory**

Check each plan against the design inventory. The final method count must be 45:

```text
Champion Mastery 4 + Champion 1 + Clash 5 + League Exp 1 + League 5 + Challenges 6 + RSO Match 3 + Status 1 + Match 4 + Spectator 1 + Summoner 2 + Tournament Stub 5 + Tournament 7 = 45
```

- [ ] **Step 2: Update README support status**

Replace the unchecked LoL line with a checked entry listing standard, Match/RSO, and Tournament clients. Do not add LoL CLI examples because CLI support is outside scope.

- [ ] **Step 3: Run formatting and full tests**

Run: `gofmt -l pkg/client pkg/lol`

Expected: no output; every implementation task formatted its files before committing.

Run: `GOTOOLCHAIN=go1.26.7 go test ./...`

Run: `GOTOOLCHAIN=go1.26.7 go test -race ./...`

Expected: all PASS.

- [ ] **Step 4: Run static and security analysis**

Run: `GOTOOLCHAIN=go1.26.7 go vet ./...`

Run: `GOTOOLCHAIN=go1.26.7 go run honnef.co/go/tools/cmd/staticcheck@latest ./...`

Run: `GOTOOLCHAIN=go1.26.7 go run golang.org/x/vuln/cmd/govulncheck@v1.7.0 ./...`

Expected: no findings and no known vulnerabilities.

- [ ] **Step 5: Verify next-Go compatibility and repository hygiene**

Run: `GOTOOLCHAIN=go1.27.0 go test ./...`

Run: `git diff --check`

Run this scan against every change since the approved design commit, including committed implementation and the working README:

```sh
if git diff 6a5c461 -- . | rg -q '^\+.*(RGAPI-[A-Za-z0-9_-]{20,}|ghp_[A-Za-z0-9]{30,}|BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY|Bearer [A-Za-z0-9._-]{20,})'; then
  echo "possible secret found"
  exit 1
fi
```

Confirm no API keys, Bearer tokens, private keys, or long credentials are present.

Expected: tests pass, diff check is clean, and no secrets are found.

- [ ] **Step 6: Commit documentation**

Run: `git add README.md && git commit -m "docs: mark LoL API clients supported"`
