package v5

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/delihiros/riot/pkg/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientPOSTEndpoints(t *testing.T) {
	emptyParticipants := []string{}
	empty := ""
	zero := 0

	tests := []struct {
		name     string
		path     string
		query    string
		body     string
		response string
		call     func(*Client) (any, error)
		want     any
	}{
		{
			name:     "create codes omits optional values and writes required zero values",
			path:     "/lol/tournament/v5/codes",
			query:    "tournamentId=0",
			body:     `{"teamSize":0,"pickType":"","mapType":"","spectatorType":"","enoughPlayers":false}`,
			response: `["code-one","code-two"]`,
			call: func(c *Client) (any, error) {
				return c.CreateTournamentCode("americas", 0, nil, TournamentCodeParametersV5{})
			},
			want: []string{"code-one", "code-two"},
		},
		{
			name:     "create codes writes explicit empty optional values and zero count",
			path:     "/lol/tournament/v5/codes",
			query:    "count=0&tournamentId=7",
			body:     `{"allowedParticipants":[],"metadata":"","teamSize":5,"pickType":"TOURNAMENT_DRAFT","mapType":"SUMMONERS_RIFT","spectatorType":"ALL","enoughPlayers":true}`,
			response: `["code/ one"]`,
			call: func(c *Client) (any, error) {
				return c.CreateTournamentCode("americas", 7, &zero, TournamentCodeParametersV5{
					AllowedParticipants: &emptyParticipants,
					Metadata:            &empty,
					TeamSize:            5,
					PickType:            "TOURNAMENT_DRAFT",
					MapType:             "SUMMONERS_RIFT",
					SpectatorType:       "ALL",
					EnoughPlayers:       true,
				})
			},
			want: []string{"code/ one"},
		},
		{
			name:     "register provider writes required empty values",
			path:     "/lol/tournament/v5/providers",
			body:     `{"region":"","url":""}`,
			response: `17`,
			call: func(c *Client) (any, error) {
				return c.RegisterProviderData("americas", ProviderRegistrationParametersV5{})
			},
			want: 17,
		},
		{
			name:     "register tournament omits name and writes required zero provider ID",
			path:     "/lol/tournament/v5/tournaments",
			body:     `{"providerId":0}`,
			response: `23`,
			call: func(c *Client) (any, error) {
				return c.RegisterTournament("americas", TournamentRegistrationParametersV5{})
			},
			want: 23,
		},
		{
			name:     "register tournament writes explicit empty name",
			path:     "/lol/tournament/v5/tournaments",
			body:     `{"providerId":9,"name":""}`,
			response: `24`,
			call: func(c *Client) (any, error) {
				return c.RegisterTournament("americas", TournamentRegistrationParametersV5{ProviderID: 9, Name: &empty})
			},
			want: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assertRequest(t, r, http.MethodPost, tt.path, tt.query, "application/json")
				assertBody(t, r, tt.body)
				return response(http.StatusOK, "200 OK", tt.response), nil
			})}

			got, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("result = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClientGETEndpoints(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(*Client) (any, error)
		want any
	}{
		{
			name: "tournament code",
			path: "/lol/tournament/v5/codes/code%2F%20one",
			body: `{"id":29,"providerId":17,"tournamentId":23,"code":"code/ one","region":"NA","map":"SUMMONERS_RIFT","teamSize":5,"spectators":"ALL","pickType":"TOURNAMENT_DRAFT","lobbyName":"Lobby 1","password":"secret","metaData":"metadata","participants":["puuid-1","puuid-2"]}`,
			call: func(c *Client) (any, error) { return c.GetTournamentCode("americas", "code/ one") },
			want: &TournamentCodeV5DTO{
				ID: 29, ProviderID: 17, TournamentID: 23, Code: "code/ one", Region: "NA",
				Map: "SUMMONERS_RIFT", TeamSize: 5, Spectators: "ALL", PickType: "TOURNAMENT_DRAFT",
				LobbyName: "Lobby 1", Password: "secret", Metadata: "metadata",
				Participants: []string{"puuid-1", "puuid-2"},
			},
		},
		{
			name: "games",
			path: "/lol/tournament/v5/games/by-code/code%2F%20one",
			body: `[{"startTime":1710000000123,"winningTeam":[{"puuid":"winner-1"},{"puuid":"winner-2"}],"losingTeam":[{"puuid":"loser-1"}],"shortCode":"code/ one","metaData":"round-1","gameId":9876543210,"gameName":"Final","gameType":"CUSTOM_GAME","gameMap":11,"gameMode":"CLASSIC","region":"NA"}]`,
			call: func(c *Client) (any, error) { return c.GetGames("americas", "code/ one") },
			want: []TournamentGamesV5{{
				StartTime:   1710000000123,
				WinningTeam: []TournamentTeamV5{{PUUID: "winner-1"}, {PUUID: "winner-2"}},
				LosingTeam:  []TournamentTeamV5{{PUUID: "loser-1"}},
				ShortCode:   "code/ one", Metadata: "round-1", GameID: 9876543210,
				GameName: "Final", GameType: "CUSTOM_GAME", GameMap: 11, GameMode: "CLASSIC", Region: "NA",
			}},
		},
		{
			name: "lobby events",
			path: "/lol/tournament/v5/lobby-events/by-code/code%2F%20one",
			body: `{"eventList":[{"timestamp":"2026-08-28T01:02:03Z","eventType":"PlayerJoined","puuid":"puuid-1"},{"timestamp":"2026-08-28T01:03:04Z","eventType":"PlayerQuit","puuid":"puuid-2"}]}`,
			call: func(c *Client) (any, error) { return c.GetLobbyEventsByCode("americas", "code/ one") },
			want: &LobbyEventV5DTOWrapper{EventList: []LobbyEventV5DTO{
				{Timestamp: "2026-08-28T01:02:03Z", EventType: "PlayerJoined", PUUID: "puuid-1"},
				{Timestamp: "2026-08-28T01:03:04Z", EventType: "PlayerQuit", PUUID: "puuid-2"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assertRequest(t, r, http.MethodGet, tt.path, "", "application/json;charset=UTF-8")
				if r.Body != nil {
					t.Fatal("GET body is non-nil")
				}
				return response(http.StatusOK, "200 OK", tt.body), nil
			})}

			got, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("result = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClientUpdateCodeBodies(t *testing.T) {
	pickType := ""
	mapType := "SUMMONERS_RIFT"
	spectatorType := "NONE"
	emptyParticipants := []string{}

	tests := []struct {
		name   string
		params *TournamentCodeUpdateParametersV5
		body   string
	}{
		{name: "nil is an empty body", params: nil, body: ""},
		{name: "zero value object", params: &TournamentCodeUpdateParametersV5{}, body: `{}`},
		{
			name: "selected fields only including empty string",
			params: &TournamentCodeUpdateParametersV5{
				PickType: &pickType, MapType: &mapType, SpectatorType: &spectatorType,
			},
			body: `{"pickType":"","mapType":"SUMMONERS_RIFT","spectatorType":"NONE"}`,
		},
		{
			name:   "explicit empty participant set",
			params: &TournamentCodeUpdateParametersV5{AllowedParticipants: &emptyParticipants},
			body:   `{"allowedParticipants":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assertRequest(t, r, http.MethodPut, "/lol/tournament/v5/codes/code%2F%20one", "", "application/json")
				assertBody(t, r, tt.body)
				return response(http.StatusNoContent, "204 No Content", ""), nil
			})}

			if err := New(client.NewWithHTTPClient("test-key", httpClient)).UpdateCode("americas", "code/ one", tt.params); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestClientReturnsNilOrZeroOnHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}
	c := New(client.NewWithHTTPClient("test-key", httpClient))

	tests := []struct {
		name string
		call func() (any, error)
		want any
	}{
		{"create tournament code", func() (any, error) { return c.CreateTournamentCode("americas", 1, nil, TournamentCodeParametersV5{}) }, []string(nil)},
		{"get tournament code", func() (any, error) { return c.GetTournamentCode("americas", "code") }, (*TournamentCodeV5DTO)(nil)},
		{"get games", func() (any, error) { return c.GetGames("americas", "code") }, []TournamentGamesV5(nil)},
		{"get lobby events", func() (any, error) { return c.GetLobbyEventsByCode("americas", "code") }, (*LobbyEventV5DTOWrapper)(nil)},
		{"register provider", func() (any, error) { return c.RegisterProviderData("americas", ProviderRegistrationParametersV5{}) }, 0},
		{"register tournament", func() (any, error) { return c.RegisterTournament("americas", TournamentRegistrationParametersV5{}) }, 0},
		{"update code", func() (any, error) { return nil, c.UpdateCode("americas", "code", nil) }, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.call()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("result = %#v, want %#v", got, tt.want)
			}
			var responseErr *client.HTTPError
			if !errors.As(err, &responseErr) {
				t.Fatalf("error = %v, want *client.HTTPError", err)
			}
			if responseErr.StatusCode != http.StatusForbidden || string(responseErr.Body) != body {
				t.Fatalf("HTTP error = %#v, want status 403 and body %q", responseErr, body)
			}
		})
	}
}

func TestClientDiscardsPartialDecodeResults(t *testing.T) {
	tests := []struct {
		name string
		body string
		call func(*Client) (any, error)
		want any
	}{
		{"create tournament code", `["decoded",1]`, func(c *Client) (any, error) {
			return c.CreateTournamentCode("americas", 1, nil, TournamentCodeParametersV5{})
		}, []string(nil)},
		{"get tournament code", `{"code":"decoded","teamSize":"bad"}`, func(c *Client) (any, error) {
			return c.GetTournamentCode("americas", "code")
		}, (*TournamentCodeV5DTO)(nil)},
		{"get games", `[{"gameName":"decoded"},{"startTime":"bad"}]`, func(c *Client) (any, error) {
			return c.GetGames("americas", "code")
		}, []TournamentGamesV5(nil)},
		{"get lobby events", `{"eventList":[{"puuid":"decoded"},{"puuid":1}]}`, func(c *Client) (any, error) {
			return c.GetLobbyEventsByCode("americas", "code")
		}, (*LobbyEventV5DTOWrapper)(nil)},
		{"register provider", `17x`, func(c *Client) (any, error) {
			return c.RegisterProviderData("americas", ProviderRegistrationParametersV5{})
		}, 0},
		{"register tournament", `{`, func(c *Client) (any, error) {
			return c.RegisterTournament("americas", TournamentRegistrationParametersV5{})
		}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, "200 OK", tt.body), nil
			})}
			got, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err == nil {
				t.Fatal("expected JSON decode error")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("result = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClientGetGamesRejectsMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetGames("americas", "code")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if got != nil {
		t.Fatalf("games = %#v, want nil", got)
	}
}

func assertRequest(t *testing.T, r *http.Request, method, path, query, contentType string) {
	t.Helper()
	if r.Method != method {
		t.Errorf("method = %q, want %q", r.Method, method)
	}
	if got := r.URL.Host; got != "americas.api.riotgames.com" {
		t.Errorf("host = %q, want americas.api.riotgames.com", got)
	}
	if got := r.URL.EscapedPath(); got != path {
		t.Errorf("escaped path = %q, want %q", got, path)
	}
	if got := r.URL.RawQuery; got != query {
		t.Errorf("query = %q, want %q", got, query)
	}
	if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
		t.Errorf("X-Riot-Token = %q, want test-key", got)
	}
	if got := r.Header.Get("Content-Type"); got != contentType {
		t.Errorf("Content-Type = %q, want %q", got, contentType)
	}
}

func assertBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(body); got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func response(statusCode int, status, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
