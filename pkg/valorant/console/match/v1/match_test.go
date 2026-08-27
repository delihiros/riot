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

func TestConsoleMatchDtoDoesNotExposePCOnlyFields(t *testing.T) {
	input := []byte(`{
		"matchInfo":{"matchId":"match-1"},
		"players":[{"puuid":"player-1"}],
		"coaches":[],
		"teams":[],
		"roundResults":[{"roundNum":1}]
	}`)
	var dto MatchDto
	if err := json.Unmarshal(input, &dto); err != nil {
		t.Fatal(err)
	}
	output, err := json.Marshal(dto)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatal(err)
	}
	matchInfo := got["matchInfo"].(map[string]interface{})
	for _, field := range []string{"gameVersion", "region", "premierMatchInfo"} {
		if _, ok := matchInfo[field]; ok {
			t.Errorf("console matchInfo unexpectedly contains %q: %s", field, output)
		}
	}
	player := got["players"].([]interface{})[0].(map[string]interface{})
	for _, field := range []string{"isObserver", "accountLevel"} {
		if _, ok := player[field]; ok {
			t.Errorf("console player unexpectedly contains %q: %s", field, output)
		}
	}
	round := got["roundResults"].([]interface{})[0].(map[string]interface{})
	if _, ok := round["winningTeamRole"]; ok {
		t.Errorf("console round unexpectedly contains winningTeamRole: %s", output)
	}
}

func TestClientGetMatchListByPUUID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		want := "https://ap.api.riotgames.com/val/match/console/v1/matchlists/by-puuid/player-1?platformType=playstation"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"puuid":"player-1","history":[{"matchId":"match-1","gameStartTimeMillis":5000000000,"queueId":"competitive"}]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchListByPUUID("ap", "player-1", "playstation")
	if err != nil {
		t.Fatal(err)
	}
	if len(dto.History) != 1 || dto.History[0].QueueID != "competitive" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}

func TestClientGetMatchByID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		want := "https://eu.api.riotgames.com/val/match/console/v1/matches/match-1"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"matchInfo":{"matchId":"match-1"},"players":[],"coaches":[],"teams":[],"roundResults":[]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchByID("eu", "match-1")
	if err != nil {
		t.Fatal(err)
	}
	if dto.MatchInfo == nil || dto.MatchInfo.MatchID != "match-1" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}

func TestClientGetRecentMatches(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		want := "https://na.api.riotgames.com/val/match/console/v1/recent-matches/by-queue/console_competitive"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"currentTime":5000000000,"matchIds":["match-1"]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetRecentMatches("na", "console_competitive")
	if err != nil {
		t.Fatal(err)
	}
	if dto.CurrentTime != 5000000000 || len(dto.MatchIDs) != 1 {
		t.Fatalf("unexpected response: %#v", dto)
	}
}
