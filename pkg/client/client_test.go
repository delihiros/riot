package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestSimpleGetRejectsInvalidRegionBeforeSendingAPIKey(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("request must not be sent")
	})}

	_, err := NewWithHTTPClient("secret", httpClient).SimpleGet("attacker.example/x", "/val/status/v1/platform-data")
	if err == nil {
		t.Fatal("expected an invalid region error")
	}
	if calls != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
	}
}

func TestSimpleGetRejectsInvalidPathBeforeSendingAPIKey(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("request must not be sent")
	})}

	for _, path := range []string{
		"https://attacker.example/x",
		"//attacker.example/x",
		"val/status/v1/platform-data",
	} {
		if _, err := NewWithHTTPClient("secret", httpClient).SimpleGet("ap", path); err == nil {
			t.Errorf("path %q: expected an invalid path error", path)
		}
	}
	if calls != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
	}
}

func TestClient_GetWithHeaderParamsReturnsErrorForNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"status":{"message":"Unauthorized","status_code":401}}`))
	}))
	defer server.Close()

	_, err := New("").GetWithHeaderParams(server.URL, nil)
	if err == nil {
		t.Fatal("expected an error for HTTP 401")
	}
}

func TestClient_GetWithHeaderParamsExposesRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "3")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	_, err := New("").GetWithHeaderParams(server.URL, nil)
	var responseErr interface {
		error
		HTTPStatusCode() int
		RetryAfterDuration() time.Duration
	}
	if !errors.As(err, &responseErr) {
		t.Fatalf("expected a typed HTTP error, got %T", err)
	}
	if responseErr.HTTPStatusCode() != http.StatusTooManyRequests {
		t.Fatalf("status code = %d, want 429", responseErr.HTTPStatusCode())
	}
	if responseErr.RetryAfterDuration() != 3*time.Second {
		t.Fatalf("retry after = %s, want 3s", responseErr.RetryAfterDuration())
	}
	if responseErr.Error() != "riot API: 429 Too Many Requests" {
		t.Fatalf("error = %q", responseErr.Error())
	}
}

func TestClientGetWithURLQueriesMergesExistingQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("existing"); got != "one" {
			t.Errorf("existing query = %q, want one", got)
		}
		if got := r.URL.Query().Get("added"); got != "two" {
			t.Errorf("added query = %q, want two", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	if _, err := New("").GetWithURLQueries(server.URL+"?existing=one#section", map[string]string{"added": "two"}); err != nil {
		t.Fatal(err)
	}
}

func TestNewSetsRequestTimeout(t *testing.T) {
	if New("").client.Timeout <= 0 {
		t.Fatal("expected a non-zero HTTP timeout")
	}
}

func TestNewWithHTTPClientPreservesProvidedSettings(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second}
	got := NewWithHTTPClient("key", httpClient).client
	if got == httpClient {
		t.Fatal("expected the provided HTTP client to be copied")
	}
	if got.Timeout != httpClient.Timeout {
		t.Fatal("provided HTTP client settings were not preserved")
	}
}

func TestNewWithNilHTTPClientUsesDefault(t *testing.T) {
	if got := NewWithHTTPClient("key", nil).client.Timeout; got != 30*time.Second {
		t.Fatalf("timeout = %s, want 30s", got)
	}
}

func TestClientDoesNotForwardAPIKeyAcrossRedirects(t *testing.T) {
	targetCalls := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetCalls++
		if got := r.Header.Get("X-Riot-Token"); got != "" {
			t.Errorf("redirected X-Riot-Token = %q, want empty", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer source.Close()

	_, err := NewWithHTTPClient("", &http.Client{}).GetWithHeaderParams(source.URL, map[string]string{
		"X-Riot-Token": "secret",
	})
	var responseErr *HTTPError
	if !errors.As(err, &responseErr) || responseErr.StatusCode != http.StatusFound {
		t.Fatalf("error = %v, want HTTP 302", err)
	}
	if targetCalls != 0 {
		t.Fatalf("redirect target calls = %d, want 0", targetCalls)
	}
}

func TestClientHaltsRequestsDuringRetryAfter(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	c := New("")
	_, _ = c.GetWithHeaderParams(server.URL, nil)
	_, err := c.GetWithHeaderParams(server.URL, nil)
	if err == nil {
		t.Fatal("expected the active rate limit to reject the request")
	}
	if calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", calls)
	}
}

func TestClient_MakeQueryString(t *testing.T) {
	type args struct {
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bar=xyz",
			args: args{
				params: map[string]string{
					"bar": "xyz",
				},
			},
			want: "?bar=xyz",
		},
		{
			name: "bar=xyz&foo=abc",
			args: args{
				params: map[string]string{
					"foo": "abc",
					"bar": "xyz",
				},
			},
			want: "?bar=xyz&foo=abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("")
			if got := c.MakeQueryString(tt.args.params); got != tt.want {
				t.Errorf("MakeQueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}
