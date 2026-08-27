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

func TestClientGetStatus(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://europe.api.riotgames.com/lor/status/v1/platform-data" {
			t.Fatalf("request URL = %s", r.URL)
		}
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Fatalf("X-Riot-Token = %q", r.Header.Get("X-Riot-Token"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"id":"europe","name":"Europe","locales":[],"maintenances":[],"incidents":[]}`)),
		}, nil
	})}

	status, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetStatus("europe")
	if err != nil {
		t.Fatal(err)
	}
	if status.ID != "europe" {
		t.Fatalf("unexpected status: %#v", status)
	}
}
