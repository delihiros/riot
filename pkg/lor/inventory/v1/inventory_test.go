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

func TestClientGetInventory(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://europe.api.riotgames.com/lor/inventory/v1/cards/me" {
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
			Body:       io.NopCloser(strings.NewReader(`[{"code":"01DE001","count":"2"}]`)),
		}, nil
	})}

	cards, err := New(client.NewWithHTTPClient("unused-api-key", httpClient)).GetInventory("europe", "Bearer access-token")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 || cards[0].Code != "01DE001" || cards[0].Count != "2" {
		t.Fatalf("unexpected cards: %#v", cards)
	}
}
