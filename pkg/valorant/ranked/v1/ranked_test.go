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

func TestClientGetLeaderboard(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		want := "https://ap.api.riotgames.com/val/ranked/v1/leaderboards/by-act/act-1?size=200&startIndex=0"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"shard":"ap","actId":"act-1","totalPlayers":1,"players":[]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeaderboard("ap", "act-1", 200, 0)
	if err != nil {
		t.Fatal(err)
	}
	if dto.ActID != "act-1" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}

func TestClientGetLeaderboardEscapesActID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		want := "https://ap.api.riotgames.com/val/ranked/v1/leaderboards/by-act/act%2Fone%3Fbad?size=200&startIndex=0"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"shard":"ap","actId":"act/one?bad","totalPlayers":0,"players":[]}`)),
		}, nil
	})}

	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeaderboard("ap", "act/one?bad", 200, 0); err != nil {
		t.Fatal(err)
	}
}

func TestLeaderboardDtoPreservesDocumentedResponse(t *testing.T) {
	input := []byte(`{
		"shard":"ap",
		"actId":"act-1",
		"totalPlayers":5000000000,
		"players":[{
			"puuid":"player-1",
			"gameName":"Player",
			"tagLine":"JP1",
			"leaderboardRank":1,
			"rankedRating":900,
			"numberOfWins":10
		}]
	}`)

	var dto LeaderboardDto
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
	if got["shard"] != "ap" || got["actId"] != "act-1" || got["totalPlayers"] != float64(5000000000) {
		t.Fatalf("leaderboard fields were not preserved: %s", output)
	}
	players, ok := got["players"].([]interface{})
	if !ok || len(players) != 1 || players[0].(map[string]interface{})["tagLine"] != "JP1" {
		t.Fatalf("players were not preserved: %s", output)
	}
}
