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
			name:     "create tournament codes without optional values",
			path:     "/lol/tournament-stub/v5/codes",
			query:    "tournamentId=-42",
			body:     `{"teamSize":0,"pickType":"","mapType":"","spectatorType":"","enoughPlayers":false}`,
			response: `["code-one","code-two"]`,
			call: func(c *Client) (any, error) {
				return c.CreateTournamentCode("americas", -42, nil, TournamentCodeParametersV5{})
			},
			want: []string{"code-one", "code-two"},
		},
		{
			name:     "create tournament codes with empty optional values and zero count",
			path:     "/lol/tournament-stub/v5/codes",
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
			name:     "register provider",
			path:     "/lol/tournament-stub/v5/providers",
			body:     `{"region":"","url":""}`,
			response: `17`,
			call: func(c *Client) (any, error) {
				return c.RegisterProviderData("americas", ProviderRegistrationParametersV5{})
			},
			want: 17,
		},
		{
			name:     "register tournament without name",
			path:     "/lol/tournament-stub/v5/tournaments",
			body:     `{"providerId":0}`,
			response: `23`,
			call: func(c *Client) (any, error) {
				return c.RegisterTournament("americas", TournamentRegistrationParametersV5{})
			},
			want: 23,
		},
		{
			name:     "register tournament with empty name",
			path:     "/lol/tournament-stub/v5/tournaments",
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
				assertRequest(t, r, http.MethodPost, tt.path, tt.query)
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %q, want application/json", got)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if got := string(body); got != tt.body {
					t.Fatalf("body = %s, want %s", got, tt.body)
				}
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
	tournamentCodeJSON := `{"code":"code/ one","lobbyName":"Lobby 1","metaData":"metadata","password":"secret","teamSize":5,"providerId":17,"pickType":"TOURNAMENT_DRAFT","tournamentId":23,"id":29,"region":"NA","map":"SUMMONERS_RIFT","participants":["puuid-1","puuid-2"]}`
	lobbyEventsJSON := `{"eventList":[{"timestamp":"2026-08-28T01:02:03Z","eventType":"PlayerJoined","puuid":"puuid-1"},{"timestamp":"2026-08-28T01:03:04Z","eventType":"PlayerQuit","puuid":"puuid-2"}]}`
	tests := []struct {
		name string
		path string
		body string
		call func(*Client) (any, error)
		want any
	}{
		{
			name: "tournament code",
			path: "/lol/tournament-stub/v5/codes/code%2F%20one",
			body: tournamentCodeJSON,
			call: func(c *Client) (any, error) { return c.GetTournamentCode("americas", "code/ one") },
			want: &TournamentCodeV5DTO{
				Code: "code/ one", LobbyName: "Lobby 1", Metadata: "metadata", Password: "secret",
				TeamSize: 5, ProviderID: 17, PickType: "TOURNAMENT_DRAFT", TournamentID: 23,
				ID: 29, Region: "NA", Map: "SUMMONERS_RIFT", Participants: []string{"puuid-1", "puuid-2"},
			},
		},
		{
			name: "lobby events",
			path: "/lol/tournament-stub/v5/lobby-events/by-code/code%2F%20one",
			body: lobbyEventsJSON,
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
				assertRequest(t, r, http.MethodGet, tt.path, "")
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
		{"get lobby events", func() (any, error) { return c.GetLobbyEventsByCode("americas", "code") }, (*LobbyEventV5DTOWrapper)(nil)},
		{"register provider", func() (any, error) { return c.RegisterProviderData("americas", ProviderRegistrationParametersV5{}) }, 0},
		{"register tournament", func() (any, error) { return c.RegisterTournament("americas", TournamentRegistrationParametersV5{}) }, 0},
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
		{"get tournament code", `{"code":"decoded","teamSize":"bad"}`, func(c *Client) (any, error) { return c.GetTournamentCode("americas", "code") }, (*TournamentCodeV5DTO)(nil)},
		{"get lobby events", `{"eventList":[{"puuid":"decoded"},{"puuid":1}]}`, func(c *Client) (any, error) { return c.GetLobbyEventsByCode("americas", "code") }, (*LobbyEventV5DTOWrapper)(nil)},
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

func assertRequest(t *testing.T, r *http.Request, method, path, query string) {
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
}

func response(statusCode int, status, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
