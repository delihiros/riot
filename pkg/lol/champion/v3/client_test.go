package v3

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

func TestClientGetChampionInfo(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		if r.URL.Path != "/lol/platform/v3/champion-rotations" {
			t.Fatalf("path = %q, want /lol/platform/v3/champion-rotations", r.URL.Path)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"newplayer":[1,2,3],"sr":[266,103,84]}`))}, nil
	})}

	info, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChampionInfo("na1")
	if err != nil {
		t.Fatal(err)
	}
	want := &ChampionInfo{NewPlayer: []int{1, 2, 3}, SR: []int{266, 103, 84}}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("info = %#v, want %#v", info, want)
	}
}

func TestClientGetChampionInfoReturnsJSONDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{"))}, nil
	})}

	info, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChampionInfo("na1")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if info != nil {
		t.Fatalf("info = %#v, want nil on decode error", info)
	}
}

func TestClientGetChampionInfoReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	info, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetChampionInfo("na1")
	if info != nil {
		t.Fatalf("info = %#v, want nil on HTTP error", info)
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
