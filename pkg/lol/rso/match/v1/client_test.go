package v1

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/delihiros/riot/pkg/client"
	matchv5 "github.com/delihiros/riot/pkg/lol/match/v5"
)

var (
	_ MatchIDsOptions = matchv5.MatchIDsOptions{}
	_ MatchDto        = matchv5.MatchDto{}
	_ TimelineDto     = matchv5.TimelineDto{}
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientRequestsAndDecodesResponses(t *testing.T) {
	const accessToken = "access-token"
	tests := []struct {
		name   string
		path   string
		body   string
		call   func(*Client) (any, error)
		assert func(*testing.T, any)
	}{
		{
			name: "match IDs",
			path: "/lol/rso-match/v1/matches/ids",
			body: `["ASIA_1","ASIA_2"]`,
			call: func(c *Client) (any, error) {
				return c.GetMatchIDs("asia", accessToken, nil)
			},
			assert: func(t *testing.T, got any) {
				t.Helper()
				if want := []string{"ASIA_1", "ASIA_2"}; !reflect.DeepEqual(got, want) {
					t.Fatalf("match IDs = %#v, want %#v", got, want)
				}
			},
		},
		{
			name: "match",
			path: "/lol/rso-match/v1/matches/ASIA_1%2F%20special",
			body: `{"metadata":{"dataVersion":"2","matchId":"ASIA_1","participants":["puuid-1"]},"info":{"gameId":123,"participants":[{"puuid":"puuid-1","summonerName":"Player"}]}}`,
			call: func(c *Client) (any, error) {
				return c.GetMatch("asia", accessToken, "ASIA_1/ special")
			},
			assert: func(t *testing.T, result any) {
				t.Helper()
				got := result.(*MatchDto)
				if got.Metadata.MatchID != "ASIA_1" || got.Info.GameID != 123 || len(got.Info.Participants) != 1 || got.Info.Participants[0].PUUID != "puuid-1" {
					t.Fatalf("match identifying values = %#v", got)
				}
			},
		},
		{
			name: "timeline",
			path: "/lol/rso-match/v1/matches/ASIA_1%2F%20special/timeline",
			body: `{"metadata":{"dataVersion":"2","matchId":"ASIA_1","participants":["puuid-1"]},"info":{"gameId":123,"frameInterval":60000,"participants":[{"participantId":1,"puuid":"puuid-1"}],"frames":[{"timestamp":60000,"events":[{"timestamp":60000,"type":"GAME_END"}]}]}}`,
			call: func(c *Client) (any, error) {
				return c.GetTimeline("asia", accessToken, "ASIA_1/ special")
			},
			assert: func(t *testing.T, result any) {
				t.Helper()
				got := result.(*TimelineDto)
				if got.Metadata.MatchID != "ASIA_1" || got.Info.GameID != 123 || len(got.Info.Participants) != 1 || got.Info.Participants[0].PUUID != "puuid-1" || len(got.Info.Frames) != 1 || got.Info.Frames[0].Events[0].Type != "GAME_END" {
					t.Fatalf("timeline identifying values = %#v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodGet {
					t.Fatalf("method = %s, want GET", r.Method)
				}
				if r.URL.Host != "asia.api.riotgames.com" {
					t.Fatalf("host = %q, want asia.api.riotgames.com", r.URL.Host)
				}
				if r.URL.EscapedPath() != tt.path {
					t.Fatalf("path = %q, want %q", r.URL.EscapedPath(), tt.path)
				}
				if r.URL.RawQuery != "" {
					t.Fatalf("query = %q, want empty", r.URL.RawQuery)
				}
				if r.Header.Get("Authorization") != "Bearer "+accessToken {
					t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
				}
				if r.Header.Get("X-Riot-Token") != "" {
					t.Fatalf("unexpected X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
				}
				if len(r.Header) != 1 {
					t.Fatalf("headers = %#v, want Authorization only", r.Header)
				}
				if strings.Contains(r.URL.String(), accessToken) {
					t.Fatalf("URL contains access token: %s", r.URL)
				}
				return testResponse(http.StatusOK, "200 OK", tt.body), nil
			})}

			got, err := tt.call(New(client.NewWithHTTPClient("sentinel-api-key", httpClient)))
			if err != nil {
				t.Fatal(err)
			}
			tt.assert(t, got)
		})
	}
}

func TestClientGetMatchIDsOptions(t *testing.T) {
	startTime := int64(1710000000)
	endTime := int64(1710003600)
	queue := 420
	matchType := "ranked & special/one"
	start := 3
	count := 100
	zero64 := int64(0)
	zero := 0
	empty := ""

	tests := []struct {
		name    string
		options *MatchIDsOptions
		query   string
	}{
		{name: "nil options"},
		{name: "nil fields", options: &MatchIDsOptions{}},
		{
			name: "full options",
			options: &MatchIDsOptions{
				StartTime: &startTime,
				EndTime:   &endTime,
				Queue:     &queue,
				Type:      &matchType,
				Start:     &start,
				Count:     &count,
			},
			query: "count=100&endTime=1710003600&queue=420&start=3&startTime=1710000000&type=ranked+%26+special%2Fone",
		},
		{
			name: "zero values",
			options: &MatchIDsOptions{
				StartTime: &zero64,
				EndTime:   &zero64,
				Queue:     &zero,
				Type:      &empty,
				Start:     &zero,
				Count:     &zero,
			},
			query: "count=0&endTime=0&queue=0&start=0&startTime=0&type=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.EscapedPath() != "/lol/rso-match/v1/matches/ids" {
					t.Fatalf("path = %q", r.URL.EscapedPath())
				}
				if r.URL.RawQuery != tt.query {
					t.Fatalf("query = %q, want %q", r.URL.RawQuery, tt.query)
				}
				return testResponse(http.StatusOK, "200 OK", `[]`), nil
			})}

			if _, err := New(client.NewWithHTTPClient("sentinel-api-key", httpClient)).GetMatchIDs("asia", "access-token", tt.options); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestClientReturnsNilOnDecodeErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		call func(*Client) (any, error)
	}{
		{
			name: "partially decoded match IDs",
			body: `["ASIA_1",2]`,
			call: func(c *Client) (any, error) { return c.GetMatchIDs("asia", "access-token", nil) },
		},
		{
			name: "partially decoded match",
			body: `{"metadata":{"matchId":"ASIA_1"},"info":{"gameId":"invalid"}}`,
			call: func(c *Client) (any, error) { return c.GetMatch("asia", "access-token", "ASIA_1") },
		},
		{
			name: "partially decoded timeline",
			body: `{"metadata":{"matchId":"ASIA_1"},"info":{"gameId":"invalid"}}`,
			call: func(c *Client) (any, error) { return c.GetTimeline("asia", "access-token", "ASIA_1") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return testResponse(http.StatusOK, "200 OK", tt.body), nil
			})}
			got, err := tt.call(New(client.NewWithHTTPClient("sentinel-api-key", httpClient)))
			if err == nil {
				t.Fatal("expected JSON decode error")
			}
			assertNilResult(t, got)
		})
	}
}

func TestClientReturnsNilOnMalformedJSON(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) (any, error)
	}{
		{name: "match IDs", call: func(c *Client) (any, error) { return c.GetMatchIDs("asia", "access-token", nil) }},
		{name: "match", call: func(c *Client) (any, error) { return c.GetMatch("asia", "access-token", "ASIA_1") }},
		{name: "timeline", call: func(c *Client) (any, error) { return c.GetTimeline("asia", "access-token", "ASIA_1") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return testResponse(http.StatusOK, "200 OK", `{`), nil
			})}
			got, err := tt.call(New(client.NewWithHTTPClient("sentinel-api-key", httpClient)))
			if err == nil {
				t.Fatal("expected malformed JSON error")
			}
			assertNilResult(t, got)
		})
	}
}

func TestClientReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	tests := []struct {
		name string
		call func(*Client) (any, error)
	}{
		{name: "match IDs", call: func(c *Client) (any, error) { return c.GetMatchIDs("asia", "access-token", nil) }},
		{name: "match", call: func(c *Client) (any, error) { return c.GetMatch("asia", "access-token", "ASIA_1") }},
		{name: "timeline", call: func(c *Client) (any, error) { return c.GetTimeline("asia", "access-token", "ASIA_1") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return testResponse(http.StatusForbidden, "403 Forbidden", body), nil
			})}
			got, err := tt.call(New(client.NewWithHTTPClient("sentinel-api-key", httpClient)))
			assertNilResult(t, got)
			var responseErr *client.HTTPError
			if !errors.As(err, &responseErr) {
				t.Fatalf("error = %v, want *client.HTTPError", err)
			}
			if responseErr.StatusCode != http.StatusForbidden || string(responseErr.Body) != body {
				t.Fatalf("HTTP error = %#v", responseErr)
			}
		})
	}
}

func TestClientRejectsInvalidRegionWithoutLeakingToken(t *testing.T) {
	const accessToken = "secret-access-token"
	tests := []struct {
		name string
		call func(*Client) (any, error)
	}{
		{name: "match IDs", call: func(c *Client) (any, error) { return c.GetMatchIDs("not/a/region", accessToken, nil) }},
		{name: "match", call: func(c *Client) (any, error) { return c.GetMatch("not/a/region", accessToken, "ASIA_1") }},
		{name: "timeline", call: func(c *Client) (any, error) { return c.GetTimeline("not/a/region", accessToken, "ASIA_1") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				return nil, errors.New("unexpected request")
			})}
			got, err := tt.call(New(client.NewWithHTTPClient("sentinel-api-key", httpClient)))
			if err == nil {
				t.Fatal("expected invalid region error")
			}
			assertNilResult(t, got)
			if calls != 0 {
				t.Fatalf("HTTP calls = %d, want 0", calls)
			}
			if strings.Contains(err.Error(), accessToken) {
				t.Fatalf("error leaks access token: %v", err)
			}
		})
	}
}

func testResponse(statusCode int, status, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func assertNilResult(t *testing.T, got any) {
	t.Helper()
	if got == nil || reflect.ValueOf(got).IsNil() {
		return
	}
	t.Fatalf("result = %#v, want nil", got)
}
