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

const challengeConfigJSON = `{"id":123456789,"localizedNames":{"en_US":{"name":"Marathon","description":"Run far"},"ja_JP":{"name":"マラソン"}},"state":"ENABLED","tracking":"LIFETIME","startTimestamp":1700000000000,"endTimestamp":1800000000000,"leaderboard":true,"thresholds":{"GOLD":10.5,"MASTER":25}}`

const playerDataJSON = `{"challenges":[{"percentile":0.75,"playersInLevel":42,"achievedTime":1712345678901,"value":27.5,"challengeId":123456789,"level":"MASTER","position":3}],"preferences":{"bannerAccent":"hextech","title":"the-title","challengeIds":["123456789","987654321"],"crestBorder":"gold","prestigeCrestBorderLevel":7},"totalPoints":{"level":"DIAMOND","current":1500,"max":2000,"precentile":91},"categoryPoints":{"TEAMWORK":{"level":"PLATINUM","current":300,"max":500,"precentile":64}}}`

func TestClientRequestsAndDecodesResponses(t *testing.T) {
	limit := 25
	tests := []struct {
		name     string
		path     string
		query    string
		response string
		want     any
		call     func(*Client) (any, error)
	}{
		{
			name:     "all challenge configs",
			path:     "/lol/challenges/v1/challenges/config",
			response: "[" + challengeConfigJSON + "]",
			want: []ChallengeConfigInfoDto{{
				ID: 123456789,
				LocalizedNames: map[string]map[string]string{
					"en_US": {"name": "Marathon", "description": "Run far"},
					"ja_JP": {"name": "マラソン"},
				},
				State: State("ENABLED"), Tracking: Tracking("LIFETIME"), StartTimestamp: 1700000000000,
				EndTimestamp: 1800000000000, Leaderboard: true, Thresholds: map[string]float64{"GOLD": 10.5, "MASTER": 25},
			}},
			call: func(c *Client) (any, error) { return c.GetAllChallengeConfigs("na1") },
		},
		{
			name:     "all challenge percentiles",
			path:     "/lol/challenges/v1/challenges/percentiles",
			response: `{"123456789":{"5":{"GOLD":0.42,"MASTER":0.01}},"987654321":{"10":{"DIAMOND":0.2}}}`,
			want: map[int64]map[int]map[string]float64{
				123456789: {5: {"GOLD": 0.42, "MASTER": 0.01}},
				987654321: {10: {"DIAMOND": 0.2}},
			},
			call: func(c *Client) (any, error) { return c.GetAllChallengePercentiles("na1") },
		},
		{
			name:     "challenge config",
			path:     "/lol/challenges/v1/challenges/123456789/config",
			response: challengeConfigJSON,
			want: &ChallengeConfigInfoDto{
				ID: 123456789,
				LocalizedNames: map[string]map[string]string{
					"en_US": {"name": "Marathon", "description": "Run far"},
					"ja_JP": {"name": "マラソン"},
				},
				State: State("ENABLED"), Tracking: Tracking("LIFETIME"), StartTimestamp: 1700000000000,
				EndTimestamp: 1800000000000, Leaderboard: true, Thresholds: map[string]float64{"GOLD": 10.5, "MASTER": 25},
			},
			call: func(c *Client) (any, error) { return c.GetChallengeConfigs("na1", 123456789) },
		},
		{
			name:     "challenge leaderboard",
			path:     "/lol/challenges/v1/challenges/123456789/leaderboards/by-level/MASTER%2FLEVEL",
			query:    "limit=25",
			response: `[{"puuid":"player/one","value":99.5,"position":1}]`,
			want:     []ApexPlayerInfoDto{{PUUID: "player/one", Value: 99.5, Position: 1}},
			call: func(c *Client) (any, error) {
				return c.GetChallengeLeaderboards("na1", 123456789, "MASTER/LEVEL", &limit)
			},
		},
		{
			name:     "challenge percentiles",
			path:     "/lol/challenges/v1/challenges/123456789/percentiles",
			response: `{"IRON":0.9,"MASTER":0.01}`,
			want:     map[string]float64{"IRON": 0.9, "MASTER": 0.01},
			call:     func(c *Client) (any, error) { return c.GetChallengePercentiles("na1", 123456789) },
		},
		{
			name:     "player data",
			path:     "/lol/challenges/v1/player-data/player%2Fone",
			response: playerDataJSON,
			want: &PlayerInfoDto{
				Challenges: []ChallengeInfoDto{{
					Percentile: 0.75, PlayersInLevel: 42, AchievedTime: 1712345678901, Value: 27.5,
					ChallengeID: 123456789, Level: "MASTER", Position: 3,
				}},
				Preferences: PlayerClientPreferencesDto{
					BannerAccent: "hextech", Title: "the-title", ChallengeIDs: []string{"123456789", "987654321"},
					CrestBorder: "gold", PrestigeCrestBorderLevel: 7,
				},
				TotalPoints: ChallengePointDto{Level: "DIAMOND", Current: 1500, Max: 2000, Precentile: 91},
				CategoryPoints: map[string]ChallengePointDto{
					"TEAMWORK": {Level: "PLATINUM", Current: 300, Max: 500, Precentile: 64},
				},
			},
			call: func(c *Client) (any, error) { return c.GetPlayerData("na1", "player/one") },
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
				if got := r.URL.EscapedPath(); got != tt.path {
					t.Fatalf("escaped path = %q, want %q", got, tt.path)
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

func TestGetAllChallengePercentilesDecodesNumericObjectKeys(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{"123456789":{"7":{"GOLD":0.5}}}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetAllChallengePercentiles("na1")
	if err != nil {
		t.Fatal(err)
	}
	levels, ok := got[int64(123456789)]
	if !ok {
		t.Fatalf("missing int64 challenge key: %#v", got)
	}
	percentiles, ok := levels[int(7)]
	if !ok {
		t.Fatalf("missing int level key: %#v", levels)
	}
	if got := percentiles["GOLD"]; got != 0.5 {
		t.Fatalf("GOLD percentile = %v, want 0.5", got)
	}
}

func TestGetChallengeLeaderboardsOmitsNilLimit(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Has("limit") {
			t.Fatalf("limit query = %q, want key absent", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", `[]`), nil
	})}

	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChallengeLeaderboards("na1", 1, "MASTER", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClientMethodsReturnNilOnMalformedJSON(t *testing.T) {
	for _, tt := range errorCalls() {
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

func TestClientMethodsReturnNilAndSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	for _, tt := range errorCalls() {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusForbidden, "403 Forbidden", body), nil
			})}

			resultIsNil, err := tt.call(New(client.NewWithHTTPClient("test-key", httpClient)))
			if !resultIsNil {
				t.Fatal("result is non-nil, want nil on HTTP error")
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
		})
	}
}

type errorCall struct {
	name string
	call func(*Client) (bool, error)
}

func errorCalls() []errorCall {
	return []errorCall{
		{"all challenge configs", func(c *Client) (bool, error) {
			got, err := c.GetAllChallengeConfigs("na1")
			return got == nil, err
		}},
		{"all challenge percentiles", func(c *Client) (bool, error) {
			got, err := c.GetAllChallengePercentiles("na1")
			return got == nil, err
		}},
		{"challenge config", func(c *Client) (bool, error) {
			got, err := c.GetChallengeConfigs("na1", 1)
			return got == nil, err
		}},
		{"challenge leaderboards", func(c *Client) (bool, error) {
			got, err := c.GetChallengeLeaderboards("na1", 1, "MASTER", nil)
			return got == nil, err
		}},
		{"challenge percentiles", func(c *Client) (bool, error) {
			got, err := c.GetChallengePercentiles("na1", 1)
			return got == nil, err
		}},
		{"player data", func(c *Client) (bool, error) {
			got, err := c.GetPlayerData("na1", "player")
			return got == nil, err
		}},
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
