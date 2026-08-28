package v4

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

const leagueListJSON = `{"leagueId":"league-id","entries":[{"freshBlood":true,"wins":12,"miniSeries":{"losses":1,"progress":"WLN","target":3,"wins":2},"inactive":true,"veteran":true,"hotStreak":true,"rank":"I","leaguePoints":88,"losses":4,"puuid":"item-player-id"}],"tier":"CHALLENGER","name":"Top Players","queue":"RANKED_SOLO_5x5"}`

const leagueEntriesJSON = `[{"leagueId":"entry-league-id","puuid":"entry-player-id","queueType":"RANKED_FLEX_SR","tier":"DIAMOND","rank":"II","leaguePoints":42,"wins":20,"losses":10,"hotStreak":true,"veteran":true,"freshBlood":true,"inactive":true,"miniSeries":{"losses":2,"progress":"LWN","target":5,"wins":3}}]`

var wantLeagueList = &LeagueListDTO{
	LeagueID: "league-id",
	Entries: []LeagueItemDTO{{
		FreshBlood:   true,
		Wins:         12,
		MiniSeries:   MiniSeriesDTO{Losses: 1, Progress: "WLN", Target: 3, Wins: 2},
		Inactive:     true,
		Veteran:      true,
		HotStreak:    true,
		Rank:         "I",
		LeaguePoints: 88,
		Losses:       4,
		PUUID:        "item-player-id",
	}},
	Tier:  "CHALLENGER",
	Name:  "Top Players",
	Queue: "RANKED_SOLO_5x5",
}

var wantLeagueEntries = []LeagueEntryDTO{{
	LeagueID:     "entry-league-id",
	PUUID:        "entry-player-id",
	QueueType:    "RANKED_FLEX_SR",
	Tier:         "DIAMOND",
	Rank:         "II",
	LeaguePoints: 42,
	Wins:         20,
	Losses:       10,
	HotStreak:    true,
	Veteran:      true,
	FreshBlood:   true,
	Inactive:     true,
	MiniSeries:   MiniSeriesDTO{Losses: 2, Progress: "LWN", Target: 5, Wins: 3},
}}

func TestClientRequestsAndDecodesResponses(t *testing.T) {
	page := 7
	tests := []struct {
		name     string
		path     string
		query    string
		response string
		want     any
		call     func(*Client) (any, error)
	}{
		{
			name:     "challenger league",
			path:     "/lol/league/v4/challengerleagues/by-queue/RANKED%2FSOLO",
			response: leagueListJSON,
			want:     wantLeagueList,
			call: func(c *Client) (any, error) {
				return c.GetChallengerLeague("na1", "RANKED/SOLO")
			},
		},
		{
			name:     "entries by PUUID",
			path:     "/lol/league/v4/entries/by-puuid/player%2Fid",
			response: leagueEntriesJSON,
			want:     wantLeagueEntries,
			call: func(c *Client) (any, error) {
				return c.GetLeagueEntriesByPUUID("na1", "player/id")
			},
		},
		{
			name:     "entries",
			path:     "/lol/league/v4/entries/FLEX%2FQUEUE/DI%2FAMOND/I%2FI",
			query:    "page=7",
			response: leagueEntriesJSON,
			want:     wantLeagueEntries,
			call: func(c *Client) (any, error) {
				return c.GetLeagueEntries("na1", "FLEX/QUEUE", "DI/AMOND", "I/I", &page)
			},
		},
		{
			name:     "grandmaster league",
			path:     "/lol/league/v4/grandmasterleagues/by-queue/GRAND%2FQUEUE",
			response: leagueListJSON,
			want:     wantLeagueList,
			call: func(c *Client) (any, error) {
				return c.GetGrandmasterLeague("na1", "GRAND/QUEUE")
			},
		},
		{
			name:     "master league",
			path:     "/lol/league/v4/masterleagues/by-queue/MASTER%2FQUEUE",
			response: leagueListJSON,
			want:     wantLeagueList,
			call: func(c *Client) (any, error) {
				return c.GetMasterLeague("na1", "MASTER/QUEUE")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodGet {
					t.Fatalf("method = %q, want GET", r.Method)
				}
				if r.URL.Host != "na1.api.riotgames.com" {
					t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
				}
				if r.URL.EscapedPath() != tt.path {
					t.Fatalf("escaped path = %q, want %q", r.URL.EscapedPath(), tt.path)
				}
				if r.URL.RawQuery != tt.query {
					t.Fatalf("query = %q, want %q", r.URL.RawQuery, tt.query)
				}
				if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
					t.Fatalf("X-Riot-Token = %q, want test-key", got)
				}
				return response(http.StatusOK, "200 OK", tt.response), nil
			})}

			got, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("response = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClientGetLeagueEntriesOmitsNilPage(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Has("page") {
			t.Fatalf("page query = %q, want key absent", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", leagueEntriesJSON), nil
	})}

	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeagueEntries("na1", "queue", "tier", "division", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClientMethodsReturnNilOnMalformedJSON(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) (bool, error)
	}{
		{"challenger league", func(c *Client) (bool, error) {
			got, err := c.GetChallengerLeague("na1", "queue")
			return got == nil, err
		}},
		{"entries by PUUID", func(c *Client) (bool, error) {
			got, err := c.GetLeagueEntriesByPUUID("na1", "puuid")
			return got == nil, err
		}},
		{"entries", func(c *Client) (bool, error) {
			got, err := c.GetLeagueEntries("na1", "queue", "tier", "division", nil)
			return got == nil, err
		}},
		{"grandmaster league", func(c *Client) (bool, error) {
			got, err := c.GetGrandmasterLeague("na1", "queue")
			return got == nil, err
		}},
		{"master league", func(c *Client) (bool, error) { got, err := c.GetMasterLeague("na1", "queue"); return got == nil, err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, "200 OK", "{"), nil
			})}

			resultIsNil, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err == nil {
				t.Fatal("expected malformed JSON error")
			}
			if !resultIsNil {
				t.Fatal("result is non-nil, want nil on decode error")
			}
		})
	}
}

func TestClientReturnsSharedHTTPErrorAndNilResult(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}

	league, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChallengerLeague("na1", "queue")
	if league != nil {
		t.Fatalf("league = %#v, want nil on HTTP error", league)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusForbidden)
	}
	if got := string(responseErr.Body); got != body {
		t.Fatalf("body = %q, want %q", got, body)
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
