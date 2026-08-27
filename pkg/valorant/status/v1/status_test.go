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

func TestClientGetStatus(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://na.api.riotgames.com/val/status/v1/platform-data" {
			t.Fatalf("request URL = %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"id":"na","name":"North America","locales":[],"maintenances":[],"incidents":[]}`)),
		}, nil
	})}

	dto, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetStatus("na")
	if err != nil {
		t.Fatal(err)
	}
	if dto.ID != "na" {
		t.Fatalf("unexpected response: %#v", dto)
	}
}

func TestPlatformDataDtoPreservesDocumentedResponse(t *testing.T) {
	input := []byte(`{
		"id":"ap",
		"name":"Asia Pacific",
		"locales":["ja-JP"],
		"maintenances":[],
		"incidents":[{
			"id":1,
			"maintenance_status":"in_progress",
			"incident_severity":"warning",
			"titles":[{"locale":"ja-JP","content":"障害"}],
			"updates":[],
			"created_at":"2026-01-01T00:00:00Z",
			"archive_at":"",
			"updated_at":"2026-01-01T00:00:00Z",
			"platforms":["windows"]
		}]
	}`)

	var dto PlatformDataDto
	if err := json.Unmarshal(input, &dto); err != nil {
		t.Fatal(err)
	}
	if dto.ID != "ap" || len(dto.Incidents) != 1 || dto.Incidents[0].IncidentSeverity != "warning" {
		t.Fatalf("unexpected status response: %#v", dto)
	}
	if _, err := json.Marshal(dto); err != nil {
		t.Fatal(err)
	}
}
