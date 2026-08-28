package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestSimpleWritesSendAuthenticatedJSONRequests(t *testing.T) {
	body := []byte{0, '"', '\n', 0xff}
	for _, tt := range []struct {
		name   string
		method string
		write  func(*Client) ([]byte, error)
	}{
		{"post", http.MethodPost, func(c *Client) ([]byte, error) {
			return c.SimplePost("na1", "/lol/example?queue=ranked&tag=a%2Fb", body)
		}},
		{"put", http.MethodPut, func(c *Client) ([]byte, error) {
			return c.SimplePut("na1", "/lol/example?queue=ranked&tag=a%2Fb", body)
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != tt.method {
					t.Errorf("method = %q, want %q", r.Method, tt.method)
				}
				if got := r.URL.String(); got != "https://na1.api.riotgames.com/lol/example?queue=ranked&tag=a%2Fb" {
					t.Errorf("URL = %q", got)
				}
				gotBody, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(gotBody, body) {
					t.Errorf("body = %v, want %v", gotBody, body)
				}
				if got := r.Header.Values("X-Riot-Token"); !reflect.DeepEqual(got, []string{"secret"}) {
					t.Errorf("X-Riot-Token = %q, want [secret]", got)
				}
				if got := r.Header.Values("Content-Type"); !reflect.DeepEqual(got, []string{"application/json"}) {
					t.Errorf("Content-Type = %q, want [application/json]", got)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(bytes.NewReader([]byte{0xfe, 0, 1})),
					Header:     make(http.Header),
				}, nil
			})}

			got, err := tt.write(NewWithHTTPClient("secret", httpClient))
			if err != nil {
				t.Fatal(err)
			}
			if want := []byte{0xfe, 0, 1}; !bytes.Equal(got, want) {
				t.Errorf("response = %v, want %v", got, want)
			}
		})
	}
}

func TestWritesWithRegionAndHeadersPreserveCallerRequest(t *testing.T) {
	body := []byte("request body")
	for _, tt := range []struct {
		name   string
		method string
		write  func(*Client, map[string]string) ([]byte, error)
	}{
		{"post", http.MethodPost, func(c *Client, headers map[string]string) ([]byte, error) {
			return c.PostWithRegionAndHeaders("euw1", "/lor/deck/v1/decks/me?existing=yes", headers, body)
		}},
		{"put", http.MethodPut, func(c *Client, headers map[string]string) ([]byte, error) {
			return c.PutWithRegionAndHeaders("euw1", "/lor/deck/v1/decks/me?existing=yes", headers, body)
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{
				"Authorization": "Bearer player",
				"Content-Type":  "application/custom+json",
				"X-Custom":      "value",
			}
			wantHeaders := map[string]string{
				"Authorization": "Bearer player",
				"Content-Type":  "application/custom+json",
				"X-Custom":      "value",
			}
			httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != tt.method {
					t.Errorf("method = %q, want %q", r.Method, tt.method)
				}
				if got := r.URL.String(); got != "https://euw1.api.riotgames.com/lor/deck/v1/decks/me?existing=yes" {
					t.Errorf("URL = %q", got)
				}
				gotBody, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(gotBody, body) {
					t.Errorf("body = %q, want %q", gotBody, body)
				}
				if len(r.Header) != len(wantHeaders) {
					t.Errorf("headers = %v, want exactly %v", r.Header, wantHeaders)
				}
				for key, want := range wantHeaders {
					if got := r.Header.Get(key); got != want {
						t.Errorf("%s = %q, want %q", key, got, want)
					}
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       http.NoBody,
					Header:     make(http.Header),
				}, nil
			})}

			if _, err := tt.write(NewWithHTTPClient("unused", httpClient), headers); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(headers, wantHeaders) {
				t.Errorf("caller headers mutated: got %v, want %v", headers, wantHeaders)
			}
		})
	}
}

func TestSimplePutNilBodySendsZeroBytes(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if len(body) != 0 || r.ContentLength != 0 {
			t.Errorf("body = %q, content length = %d; want zero bytes", body, r.ContentLength)
		}
		return &http.Response{StatusCode: http.StatusNoContent, Status: "204 No Content", Body: http.NoBody, Header: make(http.Header)}, nil
	})}

	if _, err := NewWithHTTPClient("secret", httpClient).SimplePut("na1", "/lol/example", nil); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWritesRejectInvalidRiotDestinationBeforeRoundTrip(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("request must not be sent")
	})}
	for _, tt := range []struct {
		name  string
		write func(*Client, string, string) ([]byte, error)
	}{
		{"post", func(c *Client, region, path string) ([]byte, error) { return c.SimplePost(region, path, nil) }},
		{"put", func(c *Client, region, path string) ([]byte, error) { return c.SimplePut(region, path, nil) }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := NewWithHTTPClient("secret", httpClient)
			if _, err := tt.write(c, "attacker.example/x", "/lol/example"); err == nil {
				t.Error("expected invalid routing error")
			}
			for _, path := range []string{"https://attacker.example/x", "//attacker.example/x", "lol/example"} {
				if _, err := tt.write(c, "na1", path); err == nil {
					t.Errorf("path %q: expected invalid path error", path)
				}
			}
		})
	}
	if calls != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
	}
}

func TestSimplePostRejectsRedirectWithoutForwardingCredentialsOrBody(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls > 1 {
			t.Fatalf("redirected request forwarded to %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Status:     "307 Temporary Redirect",
			Header:     http.Header{"Location": []string{"https://attacker.example/collect"}},
			Body:       http.NoBody,
			Request:    r,
		}, nil
	})}

	_, err := NewWithHTTPClient("secret", httpClient).SimplePost("na1", "/lol/example", []byte("sensitive"))
	var responseErr *HTTPError
	if !errors.As(err, &responseErr) || responseErr.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("error = %v, want HTTP 307", err)
	}
	if calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", calls)
	}
}

func TestSimpleWritesPreserveHTTPErrorAndRateLimitBehavior(t *testing.T) {
	calls := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Status:     "429 Too Many Requests",
			Header:     http.Header{"Retry-After": []string{"60"}},
			Body:       io.NopCloser(bytes.NewReader([]byte("limited"))),
		}, nil
	})}
	c := NewWithHTTPClient("secret", httpClient)

	_, err := c.SimplePost("na1", "/lol/example", nil)
	var responseErr *HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusTooManyRequests || responseErr.RetryAfter != 60*time.Second || !bytes.Equal(responseErr.Body, []byte("limited")) {
		t.Errorf("HTTP error = %+v", responseErr)
	}
	_, err = c.SimplePut("na1", "/lol/example", nil)
	if !errors.As(err, &responseErr) || responseErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("rate-limited error = %v, want HTTP 429", err)
	}
	if calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", calls)
	}
}

func TestSimplePutPreservesTransportError(t *testing.T) {
	wantErr := errors.New("round trip failed")
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, wantErr
	})}

	_, err := NewWithHTTPClient("secret", httpClient).SimplePut("na1", "/lol/example", nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapped %v", err, wantErr)
	}
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
