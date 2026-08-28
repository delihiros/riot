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

const masteryJSON = `{"puuid":"player/ one","championPointsUntilNextLevel":120,"chestGranted":true,"championId":266,"lastPlayTime":1710000000000,"championLevel":7,"championPoints":123456,"championPointsSinceLastLevel":3456,"markRequiredForNextLevel":2,"championSeasonMilestone":3,"nextSeasonMilestone":{"requireGradeCounts":{"S-":1,"S":2},"rewardMarks":4,"bonus":true,"rewardConfig":{"rewardValue":"Mythic Essence","rewardType":"CURRENCY","maximumReward":10}},"tokensEarned":5,"milestoneGrades":["S","S+"]}`

func TestClientChampionMasteryEndpoints(t *testing.T) {
	var requests []*http.Request
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests = append(requests, r)
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		body := "[" + masteryJSON + "]"
		if strings.Contains(r.URL.Path, "/by-champion/") {
			body = masteryJSON
		}
		if strings.HasSuffix(r.URL.Path, "/scores/by-puuid/player/ one") {
			body = `42`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	c := New(client.NewWithHTTPClient("test-key", httpClient))
	puuid := "player/ one"

	all, err := c.GetAllChampionMasteriesByPUUID("na1", puuid)
	if err != nil {
		t.Fatal(err)
	}
	want := ChampionMasteryDto{
		PUUID: "player/ one", ChampionPointsUntilNextLevel: 120, ChestGranted: true, ChampionID: 266,
		LastPlayTime: 1710000000000, ChampionLevel: 7, ChampionPoints: 123456, ChampionPointsSinceLastLevel: 3456,
		MarkRequiredForNextLevel: 2, ChampionSeasonMilestone: 3,
		NextSeasonMilestone: NextSeasonMilestonesDto{
			RequireGradeCounts: map[string]int{"S-": 1, "S": 2}, RewardMarks: 4, Bonus: true,
			RewardConfig: RewardConfigDto{RewardValue: "Mythic Essence", RewardType: "CURRENCY", MaximumReward: 10},
		},
		TokensEarned: 5, MilestoneGrades: []string{"S", "S+"},
	}
	if !reflect.DeepEqual(all, []ChampionMasteryDto{want}) {
		t.Fatalf("all = %#v, want %#v", all, []ChampionMasteryDto{want})
	}

	champion, err := c.GetChampionMasteryByPUUID("na1", puuid, 266)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*champion, want) {
		t.Fatalf("champion = %#v, want %#v", *champion, want)
	}

	count := 3
	top, err := c.GetTopChampionMasteriesByPUUID("na1", puuid, &count)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(top, []ChampionMasteryDto{want}) {
		t.Fatalf("top = %#v, want %#v", top, []ChampionMasteryDto{want})
	}
	zero := 0
	if _, err := c.GetTopChampionMasteriesByPUUID("na1", puuid, &zero); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetTopChampionMasteriesByPUUID("na1", puuid, nil); err != nil {
		t.Fatal(err)
	}

	score, err := c.GetChampionMasteryScoreByPUUID("na1", puuid)
	if err != nil {
		t.Fatal(err)
	}
	if score != 42 {
		t.Fatalf("score = %d, want 42", score)
	}

	wantPaths := []string{
		"/lol/champion-mastery/v4/champion-masteries/by-puuid/player%2F%20one",
		"/lol/champion-mastery/v4/champion-masteries/by-puuid/player%2F%20one/by-champion/266",
		"/lol/champion-mastery/v4/champion-masteries/by-puuid/player%2F%20one/top",
		"/lol/champion-mastery/v4/champion-masteries/by-puuid/player%2F%20one/top",
		"/lol/champion-mastery/v4/champion-masteries/by-puuid/player%2F%20one/top",
		"/lol/champion-mastery/v4/scores/by-puuid/player%2F%20one",
	}
	for i, wantPath := range wantPaths {
		if got := requests[i].URL.EscapedPath(); got != wantPath {
			t.Errorf("request %d path = %q, want %q", i, got, wantPath)
		}
	}
	if got := requests[2].URL.Query().Get("count"); got != "3" {
		t.Errorf("count query = %q, want 3", got)
	}
	if got := requests[3].URL.Query().Get("count"); got != "0" {
		t.Errorf("zero count query = %q, want 0", got)
	}
	if requests[4].URL.Query().Has("count") {
		t.Errorf("nil count query = %q, want omitted", requests[4].URL.RawQuery)
	}
}

func TestClientChampionMasteryReturnsJSONDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{"))}, nil
	})}

	mastery, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChampionMasteryByPUUID("na1", "player", 266)
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if mastery != nil {
		t.Fatalf("mastery = %#v, want nil on decode error", mastery)
	}
}

func TestClientChampionMasterySliceEndpointsReturnNilOnPartialDecodeError(t *testing.T) {
	partialJSON := "[" + masteryJSON + `,{"championId":"not-a-number"}]`
	tests := []struct {
		name string
		get  func(*Client) ([]ChampionMasteryDto, error)
	}{
		{
			name: "all",
			get: func(c *Client) ([]ChampionMasteryDto, error) {
				return c.GetAllChampionMasteriesByPUUID("na1", "player")
			},
		},
		{
			name: "top",
			get: func(c *Client) ([]ChampionMasteryDto, error) {
				return c.GetTopChampionMasteriesByPUUID("na1", "player", nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(partialJSON))}, nil
			})}

			masteries, err := tt.get(New(client.NewWithHTTPClient("test-key", httpClient)))
			if err == nil {
				t.Fatal("expected partial JSON decode error")
			}
			if masteries != nil {
				t.Fatalf("masteries = %#v, want nil on decode error", masteries)
			}
		})
	}
}

func TestClientChampionMasteryReturnsSharedHTTPError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"status":"forbidden"}`))}, nil
	})}

	_, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChampionMasteryScoreByPUUID("na1", "player")
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusForbidden)
	}
}
