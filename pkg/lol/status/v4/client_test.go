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

func TestClientGetPlatformData(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if r.URL.Host != "na1.api.riotgames.com" {
			t.Fatalf("host = %q, want na1.api.riotgames.com", r.URL.Host)
		}
		if r.URL.Path != "/lol/status/v4/platform-data" {
			t.Fatalf("path = %q, want /lol/status/v4/platform-data", r.URL.Path)
		}
		if got := r.Header.Get("X-Riot-Token"); got != "test-key" {
			t.Fatalf("X-Riot-Token = %q, want test-key", got)
		}
		return response(http.StatusOK, "200 OK", `{"id":"NA1","name":"North America","locales":["en_US","es_MX"],"maintenances":[{"id":1,"maintenance_status":"scheduled","incident_severity":"warning","titles":[{"locale":"en_US","content":"Maintenance"}],"updates":[{"id":2,"author":"Riot","publish":true,"publish_locations":["riotclient","game"],"translations":[{"locale":"en_US","content":"Starts soon"}],"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T01:00:00Z"}],"created_at":"2026-01-01T00:00:00Z","archive_at":"2026-01-02T00:00:00Z","updated_at":"2026-01-01T01:00:00Z","platforms":["NA1"]}],"incidents":[{"id":3,"maintenance_status":"in_progress","incident_severity":"critical","titles":[{"locale":"ja_JP","content":"障害"}],"updates":[{"id":4,"author":"Ops","publish":false,"publish_locations":["web"],"translations":[{"locale":"ja_JP","content":"調査中"}],"created_at":"2026-01-03T00:00:00Z","updated_at":"2026-01-03T01:00:00Z"}],"created_at":"2026-01-03T00:00:00Z","archive_at":"2026-01-04T00:00:00Z","updated_at":"2026-01-03T01:00:00Z","platforms":["NA1","PBE"]}]}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetPlatformData("na1")
	if err != nil {
		t.Fatal(err)
	}
	want := &PlatformDataDto{
		ID: "NA1", Name: "North America", Locales: []string{"en_US", "es_MX"},
		Maintenances: []StatusDto{{
			ID: 1, MaintenanceStatus: "scheduled", IncidentSeverity: "warning", Titles: []ContentDto{{Locale: "en_US", Content: "Maintenance"}},
			Updates:   []UpdateDto{{ID: 2, Author: "Riot", Publish: true, PublishLocations: []string{"riotclient", "game"}, Translations: []ContentDto{{Locale: "en_US", Content: "Starts soon"}}, CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T01:00:00Z"}},
			CreatedAt: "2026-01-01T00:00:00Z", ArchiveAt: "2026-01-02T00:00:00Z", UpdatedAt: "2026-01-01T01:00:00Z", Platforms: []string{"NA1"},
		}},
		Incidents: []StatusDto{{
			ID: 3, MaintenanceStatus: "in_progress", IncidentSeverity: "critical", Titles: []ContentDto{{Locale: "ja_JP", Content: "障害"}},
			Updates:   []UpdateDto{{ID: 4, Author: "Ops", Publish: false, PublishLocations: []string{"web"}, Translations: []ContentDto{{Locale: "ja_JP", Content: "調査中"}}, CreatedAt: "2026-01-03T00:00:00Z", UpdatedAt: "2026-01-03T01:00:00Z"}},
			CreatedAt: "2026-01-03T00:00:00Z", ArchiveAt: "2026-01-04T00:00:00Z", UpdatedAt: "2026-01-03T01:00:00Z", Platforms: []string{"NA1", "PBE"},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data = %#v, want %#v", got, want)
	}
}

func TestClientGetPlatformDataReturnsNilOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", "{"), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetPlatformData("na1")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if got != nil {
		t.Fatalf("data = %#v, want nil", got)
	}
}

func TestClientGetPlatformDataReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetPlatformData("na1")
	if got != nil {
		t.Fatalf("data = %#v, want nil", got)
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

func response(statusCode int, status, body string) *http.Response {
	return &http.Response{StatusCode: statusCode, Status: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
