package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
	RetryAfter time.Duration
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("riot API: %s", e.Status)
}

func (e *HTTPError) HTTPStatusCode() int {
	return e.StatusCode
}

func (e *HTTPError) RetryAfterDuration() time.Duration {
	return e.RetryAfter
}

type Client struct {
	apiKey          string
	client          *http.Client
	rateLimitMu     sync.Mutex
	rateLimitByHost map[string]time.Time
}

func New(apiKey string) *Client {
	return NewWithHTTPClient(apiKey, nil)
}

func NewWithHTTPClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	safeClient := *httpClient
	safeClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{
		apiKey:          apiKey,
		client:          &safeClient,
		rateLimitByHost: make(map[string]time.Time),
	}
}

func (c *Client) SimpleGet(region string, path string) ([]byte, error) {
	params := map[string]string{
		"X-Riot-Token": c.apiKey,
		"Content-Type": "application/json;charset=UTF-8",
	}
	return c.GetWithRegionAndHeaders(region, path, params)
}

func (c *Client) GetWithRegionAndHeaders(region string, path string, params map[string]string) ([]byte, error) {
	requestingURL, err := riotAPIURL(region, path)
	if err != nil {
		return nil, err
	}
	return c.GetWithHeaderParams(requestingURL, params)
}

func (c *Client) SimplePost(region, path string, body []byte) ([]byte, error) {
	return c.PostWithRegionAndHeaders(region, path, map[string]string{
		"X-Riot-Token": c.apiKey,
		"Content-Type": "application/json",
	}, body)
}

func (c *Client) PostWithRegionAndHeaders(region string, path string, params map[string]string, body []byte) ([]byte, error) {
	return c.requestWithRegionAndHeaders(http.MethodPost, region, path, params, body)
}

func (c *Client) SimplePut(region, path string, body []byte) ([]byte, error) {
	return c.PutWithRegionAndHeaders(region, path, map[string]string{
		"X-Riot-Token": c.apiKey,
		"Content-Type": "application/json",
	}, body)
}

func (c *Client) PutWithRegionAndHeaders(region, path string, params map[string]string, body []byte) ([]byte, error) {
	return c.requestWithRegionAndHeaders(http.MethodPut, region, path, params, body)
}

func (c *Client) requestWithRegionAndHeaders(method, region, path string, params map[string]string, body []byte) ([]byte, error) {
	requestingURL, err := riotAPIURL(region, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, requestingURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for key, value := range params {
		req.Header.Set(key, value)
	}
	return c.do(req)
}

func riotAPIURL(region string, path string) (string, error) {
	if region == "" || strings.IndexFunc(region, func(r rune) bool {
		return r != '-' && (r < '0' || r > '9') && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z')
	}) >= 0 {
		return "", fmt.Errorf("invalid Riot routing value %q", region)
	}
	relative, err := url.ParseRequestURI(path)
	if err != nil {
		return "", err
	}
	if relative.IsAbs() || relative.Host != "" || !strings.HasPrefix(relative.Path, "/") || strings.HasPrefix(relative.Path, "//") {
		return "", fmt.Errorf("invalid Riot API path %q", path)
	}
	relative.Scheme = "https"
	relative.Host = strings.ToLower(region) + ".api.riotgames.com"
	return relative.String(), nil
}

func (c *Client) GetWithHeaderParams(endpointURL string, params map[string]string) ([]byte, error) {
	return c.GetWithHeadersAndQueries(endpointURL, params, map[string]string{})
}

func (c *Client) GetWithURLQueries(endpointURL string, queries map[string]string) ([]byte, error) {
	return c.GetWithHeadersAndQueries(endpointURL, map[string]string{}, queries)
}

func (c *Client) GetWithHeadersAndQueries(endpointURL string, headerParams map[string]string, queries map[string]string) ([]byte, error) {
	u, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for key, value := range queries {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headerParams {
		req.Header.Set(k, v)
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	c.rateLimitMu.Lock()
	limitedUntil := c.rateLimitByHost[req.URL.Host]
	c.rateLimitMu.Unlock()
	if retryAfter := time.Until(limitedUntil); retryAfter > 0 {
		return nil, &HTTPError{
			StatusCode: http.StatusTooManyRequests,
			Status:     "429 Too Many Requests",
			RetryAfter: retryAfter,
		}
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		retryAfter, _ := strconv.Atoi(res.Header.Get("Retry-After"))
		retryAfterDuration := time.Duration(retryAfter) * time.Second
		if res.StatusCode == http.StatusTooManyRequests && retryAfterDuration > 0 {
			c.rateLimitMu.Lock()
			c.rateLimitByHost[req.URL.Host] = time.Now().Add(retryAfterDuration)
			c.rateLimitMu.Unlock()
		}
		return nil, &HTTPError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Body:       body,
			RetryAfter: retryAfterDuration,
		}
	}
	return body, nil
}

func (c *Client) MakeQueryString(params map[string]string) string {
	u := &url.URL{}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
