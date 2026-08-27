package v1

import (
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
		want := "https://ap.api.riotgames.com/val/console/ranked/v1/leaderboards/by-act/act-1?platformType=playstation&size=200&startIndex=0"
		if r.URL.String() != want {
			t.Fatalf("request URL = %s, want %s", r.URL, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"actId":"act-1",
				"totalPlayers":5000000000,
				"query":"platformType=playstation",
				"shard":"ap",
				"players":[{"tagLine":"JP1"}],
				"tierDetails":[{"tier":1}]
			}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeaderboard("ap", "act-1", "playstation", 200, 0)
	if err != nil {
		t.Fatal(err)
	}
	if dto.ActID != "act-1" || dto.TotalPlayers != 5000000000 || len(dto.Players) != 1 || dto.Players[0].TagLine != "JP1" || len(dto.TierDetails) != 1 {
		t.Fatalf("unexpected response: %#v", dto)
	}
}
