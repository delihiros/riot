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

func TestMatchDtoPreservesDocumentedResponse(t *testing.T) {
	input := []byte(`{
		"matchInfo":{
			"matchId":"match-1",
			"mapId":"map-1",
			"gameVersion":"release-1",
			"gameLengthMillis":100,
			"region":"ap",
			"gameStartMillis":5000000000,
			"provisioningFlowId":"flow-1",
			"isCompleted":true,
			"customGameName":"custom",
			"queueId":"competitive",
			"gameMode":"bomb",
			"isRanked":true,
			"seasonId":"season-1",
			"premierMatchInfo":[{"tournamentId":"tournament-1"}]
		},
		"players":[{"puuid":"player-1","isObserver":true,"accountLevel":42}],
		"coaches":[],
		"teams":[],
		"roundResults":[{
			"roundNum":1,
			"winningTeamRole":"attacker",
			"plantLocation":{"x":10,"y":20}
		}]
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
	matchInfo, ok := got["matchInfo"].(map[string]interface{})
	if !ok || matchInfo["gameVersion"] != "release-1" || matchInfo["region"] != "ap" || matchInfo["gameStartMillis"] != float64(5000000000) {
		t.Fatalf("match info was not preserved: %s", output)
	}
	if _, ok := matchInfo["premierMatchInfo"].([]interface{}); !ok {
		t.Fatalf("premier match info was not preserved: %s", output)
	}
	players := got["players"].([]interface{})
	player := players[0].(map[string]interface{})
	if player["isObserver"] != true || player["accountLevel"] != float64(42) {
		t.Fatalf("player fields were not preserved: %s", output)
	}
	rounds := got["roundResults"].([]interface{})
	round := rounds[0].(map[string]interface{})
	location := round["plantLocation"].(map[string]interface{})
	if round["winningTeamRole"] != "attacker" || location["y"] != float64(20) {
		t.Fatalf("round fields were not preserved: %s", output)
	}
}

func TestMatchListDtoPreservesQueueAndTimestamp(t *testing.T) {
	input := []byte(`{
		"puuid":"player-1",
		"history":[{
			"matchId":"match-1",
			"gameStartTimeMillis":5000000000,
			"queueId":"competitive"
		}]
	}`)

	var dto MatchListDto
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
	history, ok := got["history"].([]interface{})
	if !ok || len(history) != 1 {
		t.Fatalf("history was not preserved: %s", output)
	}
	entry := history[0].(map[string]interface{})
	if entry["matchId"] != "match-1" || entry["gameStartTimeMillis"] != float64(5000000000) || entry["queueId"] != "competitive" {
		t.Fatalf("matchlist entry was not preserved: %s", output)
	}
}

func TestClientGetMatchByID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://eu.api.riotgames.com/val/match/v1/matches/match-1" {
			t.Fatalf("request URL = %s", r.URL)
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

func TestClientGetMatchListByPUUID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://ap.api.riotgames.com/val/match/v1/matchlists/by-puuid/player-1" {
			t.Fatalf("request URL = %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"puuid":"player-1","history":[]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchListByPUUID("ap", "player-1")
	if err != nil {
		t.Fatal(err)
	}
	if dto.PUUID != "player-1" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}

func TestClientGetRecentMatches(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://ap.api.riotgames.com/val/match/v1/recent-matches/by-queue/competitive" {
			t.Fatalf("request URL = %s", r.URL)
		}
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Fatalf("X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"currentTime":5000000000,"matchIds":["match-1"]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetRecentMatches("ap", "competitive")
	if err != nil {
		t.Fatal(err)
	}
	if dto.CurrentTime != 5000000000 || len(dto.MatchIDs) != 1 || dto.MatchIDs[0] != "match-1" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}
