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

func TestClientGetMatchIDsByPUUID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/by-puuid/player%2F%20one/ids")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", `["ASIA_1","ASIA_2"]`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player/ one", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"ASIA_1", "ASIA_2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("match IDs = %#v, want %#v", got, want)
	}
}

func TestClientGetMatchIDsByPUUIDEncodesOptions(t *testing.T) {
	startTime := int64(1710000000)
	endTime := int64(1710003600)
	queue := 420
	matchType := "ranked & special/one"
	start := 3
	count := 100
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		const want = "count=100&endTime=1710003600&queue=420&start=3&startTime=1710000000&type=ranked+%26+special%2Fone"
		if r.URL.RawQuery != want {
			t.Fatalf("query = %q, want %q", r.URL.RawQuery, want)
		}
		return response(http.StatusOK, "200 OK", `[]`), nil
	})}

	options := &MatchIDsOptions{
		StartTime: &startTime,
		EndTime:   &endTime,
		Queue:     &queue,
		Type:      &matchType,
		Start:     &start,
		Count:     &count,
	}
	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", options); err != nil {
		t.Fatal(err)
	}
}

func TestClientGetMatchIDsByPUUIDIncludesZeroValues(t *testing.T) {
	zero64 := int64(0)
	zero := 0
	empty := ""
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		const want = "count=0&endTime=0&queue=0&start=0&startTime=0&type="
		if r.URL.RawQuery != want {
			t.Fatalf("query = %q, want %q", r.URL.RawQuery, want)
		}
		return response(http.StatusOK, "200 OK", `[]`), nil
	})}

	options := &MatchIDsOptions{
		StartTime: &zero64,
		EndTime:   &zero64,
		Queue:     &zero,
		Type:      &empty,
		Start:     &zero,
		Count:     &zero,
	}
	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", options); err != nil {
		t.Fatal(err)
	}
}

func TestClientGetReplay(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/by-puuid/player%2F%20one/replays")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", `{"total":2,"matchFileURLs":["https://example.test/ASIA_1.rofl","https://example.test/ASIA_2.rofl"]}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player/ one")
	if err != nil {
		t.Fatal(err)
	}
	want := &ReplayDTO{
		Total:         2,
		MatchFileURLs: []string{"https://example.test/ASIA_1.rofl", "https://example.test/ASIA_2.rofl"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("replay = %#v, want %#v", got, want)
	}
}

func TestClientGetMatchIDsByPUUIDReturnsNilOnPartialDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `["ASIA_1",2]`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", nil)
	if err == nil {
		t.Fatal("expected partial JSON decode error")
	}
	if got != nil {
		t.Fatalf("match IDs = %#v, want nil on decode error", got)
	}
}

func TestClientGetReplayReturnsNilOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if got != nil {
		t.Fatalf("replay = %#v, want nil on decode error", got)
	}
}

func TestClientGetReplayReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player")
	if got != nil {
		t.Fatalf("replay = %#v, want nil on HTTP error", got)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusForbidden)
	}
	if string(responseErr.Body) != body {
		t.Fatalf("body = %q, want %q", responseErr.Body, body)
	}
}

func assertRequest(t *testing.T, r *http.Request, path string) {
	t.Helper()
	if r.Method != http.MethodGet {
		t.Fatalf("method = %q, want GET", r.Method)
	}
	if r.URL.Host != "asia.api.riotgames.com" {
		t.Fatalf("host = %q, want asia.api.riotgames.com", r.URL.Host)
	}
	if r.URL.EscapedPath() != path {
		t.Fatalf("escaped path = %q, want %q", r.URL.EscapedPath(), path)
	}
	if r.Header.Get("X-Riot-Token") != "test-key" {
		t.Fatalf("X-Riot-Token = %q, want test-key", r.Header.Get("X-Riot-Token"))
	}
}

func response(statusCode int, status, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
