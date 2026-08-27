# League of Legends API client design

## Goal

Add typed Go clients for every League of Legends API group listed in the Riot Developer Portal on 2026-08-28. The result covers 13 API groups and 45 operations without adding CLI commands or dependencies.

## Scope

Included:

- Standard API-key endpoints
- Riot Sign On (RSO) Match endpoints
- The RSO Summoner `/me` endpoint
- Tournament v5 and Tournament Stub v5 endpoints
- All documented request and response DTOs used by those endpoints
- Tests and README support status

Excluded:

- Data Dragon
- The unsupported local League Client API
- The local Game Client API
- Teamfight Tactics APIs
- New CLI commands

The Riot Developer Portal API reference is the source of truth for paths, methods, parameters, routing families, authentication, and DTO fields.

## Package layout

Each official API group maps to one versioned package and follows the existing `client.go`, `response.go`, and `client_test.go` pattern. Packages with JSON request bodies also contain `request.go`.

| Riot API group | Go package | Operations |
| --- | --- | ---: |
| champion-mastery-v4 | `pkg/lol/championmastery/v4` | 4 |
| champion-v3 | `pkg/lol/champion/v3` | 1 |
| clash-v1 | `pkg/lol/clash/v1` | 5 |
| league-exp-v4 | `pkg/lol/leagueexp/v4` | 1 |
| league-v4 | `pkg/lol/league/v4` | 5 |
| lol-challenges-v1 | `pkg/lol/challenges/v1` | 6 |
| lol-rso-match-v1 | `pkg/lol/rso/match/v1` | 3 |
| lol-status-v4 | `pkg/lol/status/v4` | 1 |
| match-v5 | `pkg/lol/match/v5` | 4 |
| spectator-v5 | `pkg/lol/spectator/v5` | 1 |
| summoner-v4 | `pkg/lol/summoner/v4` | 2 |
| tournament-stub-v5 | `pkg/lol/tournament/stub/v5` | 5 |
| tournament-v5 | `pkg/lol/tournament/v5` | 7 |

Every package exposes `New(*client.Client) *Client`. Public method names follow the Riot operation names with Go acronym casing, such as `GetAllChampionMasteriesByPUUID` and `GetMatchIDsByPUUID`.

## Routing and authentication

Callers pass the routing value to each method:

- Platform routing, such as `jp1`, applies to Champion Mastery, Champion, Clash, League Exp, League, Challenges, Status, Spectator, and Summoner.
- Regional routing, such as `asia`, applies to Match v5 and RSO Match v1.
- Tournament callers pass the routing host appropriate for their tournament provider.

The library does not maintain a platform-to-region lookup or validate Riot enumeration values. Riot remains authoritative for supported routing, queue, tier, division, level, and tournament values.

Authentication is endpoint-specific:

- Standard and Tournament methods send the configured key only in `X-Riot-Token`.
- RSO Match methods and Summoner `GetByAccessToken` send only `Authorization: Bearer <token>`.
- Tokens never appear in URLs, query strings, errors created by this library, or logs.

## Request construction

All clients reuse `pkg/client.Client` and its HTTPS host boundary, timeout, redirect rejection, HTTP error type, and host-scoped rate-limit state.

The shared client gains `PutWithRegionAndHeaders`, implemented through the same request path as POST. It also gains `SimplePost` and `SimplePut`, which attach the configured `X-Riot-Token` and JSON content type before delegating to the header-aware methods. Tournament clients use these authenticated helpers so the private API key remains encapsulated. No generic public request builder is added.

Endpoint clients construct requests as follows:

- Escape every path parameter with `url.PathEscape`.
- Encode every query parameter with `url.Values`.
- Use direct method parameters for required values and endpoints with at most two optional values.
- Use endpoint-specific `Options` structs when an endpoint has several optional filters.
- Represent optional numeric values with pointers when zero differs from omission.
- Encode Tournament request DTOs with `encoding/json` and set `Content-Type: application/json`.

Required Tournament request fields use value types and have no `omitempty` tag. Every optional field uses a pointer and `omitempty`; optional collections use pointers to slices so callers can distinguish an omitted collection from an explicitly empty JSON array. All fields in `TournamentCodeUpdateParametersV5` are optional. `UpdateCode` therefore accepts a pointer to that DTO: a nil pointer sends the officially permitted empty body, while a non-nil pointer encodes only the selected updates.

## Response model

Response structs reproduce every field and JSON name documented by Riot for the operation. Go types map as follows:

| Riot type | Go type |
| --- | --- |
| string | `string` |
| boolean | `bool` |
| int | `int` |
| long | `int64` |
| float or double | `float64` |
| List[T] | `[]T` |
| Map[string, T] | `map[string]T` |
| Set[T] | `[]T` |
| Map[Long, T] | `map[int64]T` |
| Map[Integer, T] | `map[int]T` |
| Map[Level, T] | `map[string]T` |
| documented DTO | named Go struct |

JSON represents Riot sets as arrays, so the library preserves wire order and does not add set semantics. JSON object keys are decoded into the declared numeric or string map key type. Object responses return pointers, list and set responses return slices, and primitive responses return their matching Go values. JSON decoding errors are returned to the caller. Transport and non-2xx responses keep the existing shared-client behavior.

## Endpoint inventory

### Champion Mastery v4

- `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}`
- `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}/by-champion/{championId}`
- `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}/top`
- `GET /lol/champion-mastery/v4/scores/by-puuid/{encryptedPUUID}`

### Champion v3

- `GET /lol/platform/v3/champion-rotations`

### Clash v1

- `GET /lol/clash/v1/players/by-puuid/{puuid}`
- `GET /lol/clash/v1/teams/{teamId}`
- `GET /lol/clash/v1/tournaments`
- `GET /lol/clash/v1/tournaments/by-team/{teamId}`
- `GET /lol/clash/v1/tournaments/{tournamentId}`

### League Exp v4

- `GET /lol/league-exp/v4/entries/{queue}/{tier}/{division}`

### League v4

- `GET /lol/league/v4/challengerleagues/by-queue/{queue}`
- `GET /lol/league/v4/entries/by-puuid/{encryptedPUUID}`
- `GET /lol/league/v4/entries/{queue}/{tier}/{division}`
- `GET /lol/league/v4/grandmasterleagues/by-queue/{queue}`
- `GET /lol/league/v4/masterleagues/by-queue/{queue}`

### Challenges v1

- `GET /lol/challenges/v1/challenges/config`
- `GET /lol/challenges/v1/challenges/percentiles`
- `GET /lol/challenges/v1/challenges/{challengeId}/config`
- `GET /lol/challenges/v1/challenges/{challengeId}/leaderboards/by-level/{level}`
- `GET /lol/challenges/v1/challenges/{challengeId}/percentiles`
- `GET /lol/challenges/v1/player-data/{puuid}`

### RSO Match v1

- `GET /lol/rso-match/v1/matches/ids`
- `GET /lol/rso-match/v1/matches/{matchId}`
- `GET /lol/rso-match/v1/matches/{matchId}/timeline`

### Status v4

- `GET /lol/status/v4/platform-data`

### Match v5

- `GET /lol/match/v5/matches/by-puuid/{puuid}/ids`
- `GET /lol/match/v5/matches/by-puuid/{puuid}/replays`
- `GET /lol/match/v5/matches/{matchId}`
- `GET /lol/match/v5/matches/{matchId}/timeline`

### Spectator v5

- `GET /lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}`

### Summoner v4

- `GET /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}`
- `GET /lol/summoner/v4/summoners/me`

### Tournament Stub v5

- `POST /lol/tournament-stub/v5/codes`
- `GET /lol/tournament-stub/v5/codes/{tournamentCode}`
- `GET /lol/tournament-stub/v5/lobby-events/by-code/{tournamentCode}`
- `POST /lol/tournament-stub/v5/providers`
- `POST /lol/tournament-stub/v5/tournaments`

### Tournament v5

- `POST /lol/tournament/v5/codes`
- `GET /lol/tournament/v5/codes/{tournamentCode}`
- `PUT /lol/tournament/v5/codes/{tournamentCode}`
- `GET /lol/tournament/v5/games/by-code/{tournamentCode}`
- `GET /lol/tournament/v5/lobby-events/by-code/{tournamentCode}`
- `POST /lol/tournament/v5/providers`
- `POST /lol/tournament/v5/tournaments`

## Test strategy

Each production method is introduced through a red-green-refactor cycle. Package tests use an injected `http.RoundTripper` at the network boundary and assert observable request and response behavior:

- HTTP method and routing host
- Escaped path and encoded query parameters
- Correct API-key or Bearer header, with the other credential absent
- Complete documented JSON request bodies
- Complete representative response DTO decoding
- Primitive and list response decoding
- JSON errors and shared HTTP errors

The shared client receives a regression test for PUT requests. Existing tests continue to cover invalid routing, invalid paths, secret-safe redirects, timeouts, non-2xx responses, and 429 handling.

Final verification runs:

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `staticcheck ./...`
- `govulncheck ./...`
- `go test ./...` with Go 1.27
- `git diff --check`

## Delivery

Implementation may be organized into standard platform APIs, Match and RSO APIs, and Tournament APIs, but the LoL support checkbox is marked complete only after all 45 operations and their documented DTOs are present and verified.
