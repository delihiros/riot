package v1

import (
	"reflect"
	"riot/config"
	"riot/pkg/client"
	"testing"
)

func TestClient_GetMatchByID(t *testing.T) {
	config.ReadConfig("")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		region  string
		matchID string
	}
	var tests []struct {
		name    string
		args    args
		want    *MatchDto
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetMatchByID(tt.args.region, tt.args.matchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
