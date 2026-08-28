package v5

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

func TestClientGetCurrentGameInfoByPUUID(t *testing.T) {
	const body = `{"gameId":9007199254740991,"gameType":"MATCHED_GAME","gameStartTime":1710000000000,"mapId":4294967296,"gameLength":1234,"platformId":"NA1","gameMode":"CLASSIC","bannedChampions":[{"pickTurn":1,"championId":4294967297,"teamId":4294967298}],"gameQueueConfigId":420,"observers":{"encryptionKey":"observer-key"},"participants":[{"championId":4294967299,"perks":{"perkIds":[4294967300,4294967301],"perkStyle":4294967302,"perkSubStyle":4294967303},"profileIconId":4294967304,"bot":true,"teamId":4294967305,"puuid":"player/ one","spell1Id":4294967306,"spell2Id":4294967307,"gameCustomizationObjects":[{"category":"summonerEmote","content":"emote-1"}]}]}`
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		const wantPath = "/lol/spectator/v5/active-games/by-summoner/player%2F%20one"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Fatalf("path = %q, want %q", got, wantPath)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}

	game, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetCurrentGameInfoByPUUID("na1", "player/ one")
	if err != nil {
		t.Fatal(err)
	}
	want := &CurrentGameInfo{
		GameID:            9007199254740991,
		GameType:          "MATCHED_GAME",
		GameStartTime:     1710000000000,
		MapID:             4294967296,
		GameLength:        1234,
		PlatformID:        "NA1",
		GameMode:          "CLASSIC",
		BannedChampions:   []BannedChampion{{PickTurn: 1, ChampionID: 4294967297, TeamID: 4294967298}},
		GameQueueConfigID: 420,
		Observers:         Observer{EncryptionKey: "observer-key"},
		Participants: []CurrentGameParticipant{{
			ChampionID:    4294967299,
			Perks:         Perks{PerkIDs: []int64{4294967300, 4294967301}, PerkStyle: 4294967302, PerkSubStyle: 4294967303},
			ProfileIconID: 4294967304,
			Bot:           true,
			TeamID:        4294967305,
			PUUID:         "player/ one",
			Spell1ID:      4294967306,
			Spell2ID:      4294967307,
			GameCustomizationObjects: []GameCustomizationObject{{
				Category: "summonerEmote",
				Content:  "emote-1",
			}},
		}},
	}
	if !reflect.DeepEqual(game, want) {
		t.Fatalf("game = %#v, want %#v", game, want)
	}
}

func TestClientGetCurrentGameInfoByPUUIDReturnsNilOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("{")),
		}, nil
	})}

	game, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetCurrentGameInfoByPUUID("na1", "player")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if game != nil {
		t.Fatalf("game = %#v, want nil on decode error", game)
	}
}

func TestClientGetCurrentGameInfoByPUUIDReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"not found"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}

	game, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetCurrentGameInfoByPUUID("na1", "player")
	if game != nil {
		t.Fatalf("game = %#v, want nil on HTTP error", game)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusNotFound)
	}
	if string(responseErr.Body) != body {
		t.Fatalf("body = %q, want %q", responseErr.Body, body)
	}
}
