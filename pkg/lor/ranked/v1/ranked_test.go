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

func TestClientGetLeaderboards(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://americas.api.riotgames.com/lor/ranked/v1/leaderboards" {
			t.Fatalf("request URL = %s", r.URL)
		}
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Fatalf("X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"players":[{"name":"Player","rank":1,"lp":100}]}`)),
		}, nil
	})}

	leaderboard, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetLeaderboards("americas")
	if err != nil {
		t.Fatal(err)
	}
	if len(leaderboard.Players) != 1 || leaderboard.Players[0].LP != 100 {
		t.Fatalf("unexpected leaderboard: %#v", leaderboard)
	}
}
