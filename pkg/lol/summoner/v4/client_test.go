package v4

import (
	"errors"
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

func TestClientGetByPUUID(t *testing.T) {
	puuid := "player/ one"
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		if got, want := r.URL.EscapedPath(), "/lol/summoner/v4/summoners/by-puuid/player%2F%20one"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"profileIconId":12,"revisionDate":1710000000000,"puuid":"player/ one","summonerLevel":42}`)),
		}, nil
	})}

	summoner, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetByPUUID("na1", puuid)
	if err != nil {
		t.Fatal(err)
	}
	want := &SummonerDTO{ProfileIconID: 12, RevisionDate: 1710000000000, PUUID: puuid, SummonerLevel: 42}
	if *summoner != *want {
		t.Fatalf("summoner = %#v, want %#v", summoner, want)
	}
}

func TestClientGetByPUUIDReturnsNilOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{"))}, nil
	})}

	summoner, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetByPUUID("na1", "player")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if summoner != nil {
		t.Fatalf("summoner = %#v, want nil", summoner)
	}
}

func TestClientGetByPUUIDReturnsSharedHTTPError(t *testing.T) {
	body := `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	summoner, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetByPUUID("na1", "player")
	if summoner != nil {
		t.Fatalf("summoner = %#v, want nil", summoner)
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
