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

func TestClientGetUserDeck(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://americas.api.riotgames.com/lor/deck/v1/decks/me" {
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
			Body:       io.NopCloser(strings.NewReader(`[{"id":"deck-1","name":"Deck","code":"CODE"}]`)),
		}, nil
	})}

	decks, err := New(client.NewWithHTTPClient("unused-api-key", httpClient)).GetUserDeck("americas", "Bearer access-token")
	if err != nil {
		t.Fatal(err)
	}
	if len(decks) != 1 || decks[0].ID != "deck-1" || decks[0].Code != "CODE" {
		t.Fatalf("unexpected decks: %#v", decks)
	}
}

func TestClientCreateDeck(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("request method = %s", r.Method)
		}
		if r.URL.String() != "https://americas.api.riotgames.com/lor/deck/v1/decks/me" {
			t.Fatalf("request URL = %s", r.URL)
		}
		if r.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		var deck NewDeckDto
		if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
			t.Fatal(err)
		}
		if deck.Name != "Deck" || deck.Code != "CODE" {
			t.Fatalf("unexpected request body: %#v", deck)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`"deck-1"`)),
		}, nil
	})}

	id, err := New(client.NewWithHTTPClient("unused-api-key", httpClient)).CreateDeck(
		"americas",
		"Bearer access-token",
		&NewDeckDto{Name: "Deck", Code: "CODE"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if id != "deck-1" {
		t.Fatalf("deck ID = %q", id)
	}
}
