package v1

import (
	"reflect"
	"riot/config"
	"riot/pkg/client"
	"testing"
)

func TestClient_GetContent(t *testing.T) {
	config.ReadConfig("")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "version check",
			want:    "release-02.09",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetContent("ap", "ja-JP")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Version, tt.want) {
				t.Errorf("GetContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
