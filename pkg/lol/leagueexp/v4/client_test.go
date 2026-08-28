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

const leagueEntriesJSON = `[{"leagueId":"league-id","summonerId":"summoner-id","puuid":"player-id","queueType":"RANKED_SOLO_5x5","tier":"GOLD","rank":"II","leaguePoints":42,"wins":20,"losses":10,"hotStreak":true,"veteran":true,"freshBlood":true,"inactive":true,"miniSeries":{"losses":1,"progress":"WLN","target":3,"wins":2}}]`

func TestClientGetLeagueEntriesEscapesPathAndEncodesPage(t *testing.T) {
	var request *http.Request
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		request = r
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(leagueEntriesJSON))}, nil
	})}

	page := 7
	entries, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeagueEntries("na1", "RANKED/SOLO", "GO/LD", "I/I", &page)
	if err != nil {
		t.Fatal(err)
	}
	want := LeagueEntryDTO{
		LeagueID: "league-id", SummonerID: "summoner-id", PUUID: "player-id", QueueType: "RANKED_SOLO_5x5", Tier: "GOLD", Rank: "II",
		LeaguePoints: 42, Wins: 20, Losses: 10, HotStreak: true, Veteran: true, FreshBlood: true, Inactive: true,
		MiniSeries: MiniSeriesDTO{Losses: 1, Progress: "WLN", Target: 3, Wins: 2},
	}
	if !reflect.DeepEqual(entries, []LeagueEntryDTO{want}) {
		t.Fatalf("entries = %#v, want %#v", entries, []LeagueEntryDTO{want})
	}
	if got, want := request.URL.EscapedPath(), "/lol/league-exp/v4/entries/RANKED%2FSOLO/GO%2FLD/I%2FI"; got != want {
		t.Errorf("EscapedPath() = %q, want %q", got, want)
	}
	if got, want := request.URL.RawPath, "/lol/league-exp/v4/entries/RANKED%2FSOLO/GO%2FLD/I%2FI"; got != want {
		t.Errorf("RawPath = %q, want %q", got, want)
	}
	if got := request.URL.Query().Get("page"); got != "7" {
		t.Errorf("page = %q, want 7", got)
	}
}

func TestClientGetLeagueEntriesOmitsNilPage(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Has("page") {
			t.Fatalf("page query = %q, want omitted", r.URL.RawQuery)
		}
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(leagueEntriesJSON))}, nil
	})}

	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeagueEntries("na1", "RANKED_SOLO_5x5", "GOLD", "II", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClientGetLeagueEntriesReturnsJSONDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{"))}, nil
	})}

	entries, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeagueEntries("na1", "queue", "tier", "division", nil)
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if entries != nil {
		t.Fatalf("entries = %#v, want nil on decode error", entries)
	}
}

func TestClientGetLeagueEntriesReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	entries, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeagueEntries("na1", "queue", "tier", "division", nil)
	if entries != nil {
		t.Fatalf("entries = %#v, want nil on HTTP error", entries)
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
