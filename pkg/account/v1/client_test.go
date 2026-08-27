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

func TestClientGetAccountByAccessToken(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://asia.api.riotgames.com/riot/account/v1/accounts/me" {
			t.Fatalf("request URL = %s", r.URL)
		}
		if r.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Riot-Token") != "" {
			t.Fatalf("unexpected X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"puuid":"player-1","gameName":"Player","tagLine":"JP1"}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("", httpClient)).GetAccountByAccessToken("asia", "access-token")
	if err != nil {
		t.Fatal(err)
	}
	want := AccountDto{PUUID: "player-1", GameName: "Player", TagLine: "JP1"}
	got, _ := json.Marshal(dto)
	expected, _ := json.Marshal(want)
	if string(got) != string(expected) {
		t.Fatalf("account = %s, want %s", got, expected)
	}
}

func TestClientGetAccountByAccessTokenRejectsInvalidRegionBeforeSendingToken(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, nil
	})}

	_, err := New(client.NewWithHTTPClient("", httpClient)).GetAccountByAccessToken("attacker.example/x", "access-token")
	if err == nil {
		t.Fatal("expected an invalid region error")
	}
	if calls != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
	}
}

func TestClientAPIKeyEndpoints(t *testing.T) {
	var paths []string
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		paths = append(paths, r.URL.EscapedPath())
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Fatalf("X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		body := `{"puuid":"player-1","gameName":"Player Name","tagLine":"JP1"}`
		if strings.Contains(r.URL.Path, "active-shards") {
			body = `{"puuid":"player-1","game":"val","activeShard":"ap"}`
		}
		if strings.Contains(r.URL.Path, "/region/") {
			body = `{"puuid":"player-1","game":"lol","region":"na1"}`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	c := New(client.NewWithHTTPClient("test-key", httpClient))

	if _, err := c.GetAccountByPUUID("asia", "player-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetAccountByRiotID("asia", "Player Name", "JP1"); err != nil {
		t.Fatal(err)
	}
	shard, err := c.GetActiveShard("asia", "val", "player-1")
	if err != nil {
		t.Fatal(err)
	}
	if shard.ActiveShard != "ap" {
		t.Fatalf("active shard = %q", shard.ActiveShard)
	}
	accountRegion, err := c.GetActiveRegion("americas", "lol", "player-1")
	if err != nil {
		t.Fatal(err)
	}
	if accountRegion.Region != "na1" {
		t.Fatalf("account region = %q", accountRegion.Region)
	}

	want := []string{
		"/riot/account/v1/accounts/by-puuid/player-1",
		"/riot/account/v1/accounts/by-riot-id/Player%20Name/JP1",
		"/riot/account/v1/active-shards/by-game/val/by-puuid/player-1",
		"/riot/account/v1/region/by-game/lol/by-puuid/player-1",
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("path %d = %q, want %q", i, paths[i], want[i])
		}
	}
}
