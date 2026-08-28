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
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("unexpected Authorization = %q", got)
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

func TestClientGetByAccessToken(t *testing.T) {
	const accessToken = "access-token"
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Host, "na1.api.riotgames.com"; got != want {
			t.Fatalf("host = %q, want %q", got, want)
		}
		if got, want := r.URL.EscapedPath(), "/lol/summoner/v4/summoners/me"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got := r.URL.RawQuery; got != "" {
			t.Fatalf("query = %q, want empty", got)
		}
		if got, want := r.Header.Get("Authorization"), "Bearer "+accessToken; got != want {
			t.Fatalf("Authorization = %q, want %q", got, want)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "" {
			t.Fatalf("unexpected X-Riot-Token = %q", got)
		}
		if len(r.Header) != 1 {
			t.Fatalf("headers = %#v, want Authorization only", r.Header)
		}
		if strings.Contains(r.URL.String(), accessToken) {
			t.Fatalf("URL contains access token: %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"profileIconId":12,"revisionDate":1710000000000,"puuid":"player-1","summonerLevel":42}`)),
		}, nil
	})}

	summoner, err := New(client.NewWithHTTPClient("sentinel-api-key", httpClient)).GetByAccessToken("na1", accessToken)
	if err != nil {
		t.Fatal(err)
	}
	want := &SummonerDTO{ProfileIconID: 12, RevisionDate: 1710000000000, PUUID: "player-1", SummonerLevel: 42}
	if *summoner != *want {
		t.Fatalf("summoner = %#v, want %#v", summoner, want)
	}
}

func TestClientGetByAccessTokenReturnsNilOnErrors(t *testing.T) {
	tests := []struct {
		name      string
		transport roundTripFunc
		assert    func(*testing.T, error)
	}{
		{
			name: "malformed JSON",
			transport: func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{"))}, nil
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected malformed JSON error")
				}
			},
		},
		{
			name: "partially decoded JSON",
			transport: func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"profileIconId":12,"revisionDate":"invalid"}`))}, nil
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected JSON error")
				}
			},
		},
		{
			name: "request error",
			transport: func(*http.Request) (*http.Response, error) {
				return nil, errors.New("request failed")
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected request error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summoner, err := New(client.NewWithHTTPClient("sentinel-api-key", &http.Client{Transport: tt.transport})).GetByAccessToken("na1", "access-token")
			if summoner != nil {
				t.Fatalf("summoner = %#v, want nil", summoner)
			}
			tt.assert(t, err)
		})
	}
}

func TestClientGetByAccessTokenReturnsSharedHTTPError(t *testing.T) {
	body := `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	summoner, err := New(client.NewWithHTTPClient("sentinel-api-key", httpClient)).GetByAccessToken("na1", "access-token")
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

func TestClientGetByAccessTokenDoesNotExposeTokenOnInvalidRegion(t *testing.T) {
	const accessToken = "secret-access-token"
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("request must not be sent")
	})}

	summoner, err := New(client.NewWithHTTPClient("sentinel-api-key", httpClient)).GetByAccessToken("attacker.example/x", accessToken)
	if summoner != nil {
		t.Fatalf("summoner = %#v, want nil", summoner)
	}
	if err == nil {
		t.Fatal("expected invalid region error")
	}
	if strings.Contains(err.Error(), accessToken) {
		t.Fatalf("error exposes access token: %v", err)
	}
	if calls != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
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
