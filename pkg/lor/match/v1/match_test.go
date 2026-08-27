package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/delihiros/riot/pkg/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientMatchEndpoints(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Fatalf("X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		var body string
		switch r.URL.String() {
		case "https://sea.api.riotgames.com/lor/match/v1/matches/by-puuid/player-1/ids":
			body = `["match-1"]`
		case "https://sea.api.riotgames.com/lor/match/v1/matches/match-1":
			body = `{"metadata":{"match_id":"match-1"},"info":{"game_format":"standard"}}`
		default:
			t.Fatalf("unexpected request URL: %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	c := New(client.NewWithHTTPClient("test-key", httpClient))

	ids, err := c.GetMatchListByPUUID("sea", "player-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "match-1" {
		t.Fatalf("unexpected match IDs: %#v", ids)
	}
	match, err := c.GetMatchByID("sea", "match-1")
	if err != nil {
		t.Fatal(err)
	}
	output, err := json.Marshal(match)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatal(err)
	}
	info, ok := got["info"].(map[string]interface{})
	if !ok || info["game_format"] != "standard" {
		t.Fatalf("match info was not preserved: %s", output)
	}
}
