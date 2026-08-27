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

func TestClientGetContentOmitsEmptyLocale(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://ap.api.riotgames.com/val/content/v1/contents" {
			t.Fatalf("request URL = %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"version":"release-current"}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetContent("ap", "")
	if err != nil {
		t.Fatal(err)
	}
	if dto.Version != "release-current" {
		t.Fatalf("version = %q", dto.Version)
	}
}

func TestClientGetContentIncludesLocale(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://eu.api.riotgames.com/val/content/v1/contents?locale=fr-FR" {
			t.Fatalf("request URL = %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"version":"release-current"}`)),
		}, nil
	})}

	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetContent("eu", "fr-FR"); err != nil {
		t.Fatal(err)
	}
}

func TestContentDtoPreservesSprayLevelsAndActs(t *testing.T) {
	input := []byte(`{
		"version":"release-current",
		"sprayLevels":[{"name":"Level 1","id":"spray-level-1"}],
		"acts":[{
			"name":"Act 1",
			"localizedNames":{"ja-JP":"アクト1"},
			"id":"act-1",
			"isActive":true
		}]
	}`)

	var dto ContentDto
	if err := json.Unmarshal(input, &dto); err != nil {
		t.Fatal(err)
	}
	output, err := json.Marshal(dto)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatal(err)
	}
	if levels, ok := got["sprayLevels"].([]interface{}); !ok || len(levels) != 1 {
		t.Fatalf("spray levels were not preserved: %s", output)
	}
	acts, ok := got["acts"].([]interface{})
	if !ok || len(acts) != 1 {
		t.Fatalf("acts were not preserved: %s", output)
	}
	act := acts[0].(map[string]interface{})
	if act["id"] != "act-1" || act["isActive"] != true {
		t.Fatalf("act fields were not preserved: %s", output)
	}
}
