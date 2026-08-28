package v1

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

const (
	playerJSON     = `{"puuid":"player/one","position":"MIDDLE","role":"CAPTAIN"}`
	teamJSON       = `{"id":"team/one","tournamentId":42,"name":"Example Team","iconId":7,"tier":3,"captain":"player/one","abbreviation":"EX","players":[` + playerJSON + `]}`
	tournamentJSON = `{"id":42,"themeId":9,"nameKey":"name.primary","nameKeySecondary":"name.secondary","schedule":[{"id":5,"registrationTime":1710000000000,"startTime":1710100000000,"cancelled":true}]}`
)

func TestClientEndpoints(t *testing.T) {
	wantPlayer := PlayerDto{PUUID: "player/one", Position: "MIDDLE", Role: "CAPTAIN"}
	wantTeam := &TeamDto{
		ID: "team/one", TournamentID: 42, Name: "Example Team", IconID: 7, Tier: 3,
		Captain: "player/one", Abbreviation: "EX", Players: []PlayerDto{wantPlayer},
	}
	wantTournament := &TournamentDto{
		ID: 42, ThemeID: 9, NameKey: "name.primary", NameKeySecondary: "name.secondary",
		Schedule: []TournamentPhaseDto{{ID: 5, RegistrationTime: 1710000000000, StartTime: 1710100000000, Cancelled: true}},
	}

	tests := []struct {
		name string
		path string
		body string
		call func(*Client) (any, error)
		want any
	}{
		{
			name: "players by PUUID",
			path: "/lol/clash/v1/players/by-puuid/player%2Fone",
			body: "[" + playerJSON + "]",
			call: func(c *Client) (any, error) { return c.GetPlayersByPUUID("euw1", "player/one") },
			want: []PlayerDto{wantPlayer},
		},
		{
			name: "team by ID",
			path: "/lol/clash/v1/teams/team%2Fone",
			body: teamJSON,
			call: func(c *Client) (any, error) { return c.GetTeamByID("euw1", "team/one") },
			want: wantTeam,
		},
		{
			name: "tournaments",
			path: "/lol/clash/v1/tournaments",
			body: "[" + tournamentJSON + "]",
			call: func(c *Client) (any, error) { return c.GetTournaments("euw1") },
			want: []TournamentDto{*wantTournament},
		},
		{
			name: "tournament by team",
			path: "/lol/clash/v1/tournaments/by-team/team%2Fone",
			body: tournamentJSON,
			call: func(c *Client) (any, error) { return c.GetTournamentByTeam("euw1", "team/one") },
			want: wantTournament,
		},
		{
			name: "tournament by ID",
			path: "/lol/clash/v1/tournaments/42",
			body: tournamentJSON,
			call: func(c *Client) (any, error) { return c.GetTournamentByID("euw1", 42) },
			want: wantTournament,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %q, want GET", r.Method)
				}
				if r.URL.Host != "euw1.api.riotgames.com" {
					t.Errorf("host = %q, want euw1.api.riotgames.com", r.URL.Host)
				}
				if got := r.URL.EscapedPath(); got != tt.path {
					t.Errorf("escaped path = %q, want %q", got, tt.path)
				}
				if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
					t.Errorf("X-Riot-Token = %q, want test-key", got)
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

func TestClientReturnsNilResultsOnHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}
	c := New(client.NewWithHTTPClient("test-key", httpClient))

	calls := []struct {
		name string
		call func() (any, error)
	}{
		{"players by PUUID", func() (any, error) { return c.GetPlayersByPUUID("euw1", "player") }},
		{"team by ID", func() (any, error) { return c.GetTeamByID("euw1", "team") }},
		{"tournaments", func() (any, error) { return c.GetTournaments("euw1") }},
		{"tournament by team", func() (any, error) { return c.GetTournamentByTeam("euw1", "team") }},
		{"tournament by ID", func() (any, error) { return c.GetTournamentByID("euw1", 42) }},
	}

	for _, tt := range calls {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.call()
			if got != nil && !reflect.ValueOf(got).IsNil() {
				t.Fatalf("result = %#v, want nil", got)
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

func TestClientReturnsNilResultOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{"id":42,"players":`), nil
	})}

	team, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetTeamByID("euw1", "team")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if team != nil {
		t.Fatalf("team = %#v, want nil", team)
	}
}

func response(statusCode int, status string, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
