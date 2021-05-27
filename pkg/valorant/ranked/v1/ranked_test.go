package v1

import (
	"reflect"
	"riot/config"
	"riot/pkg/client"
	"testing"
)

func TestClient_GetLeaderboard(t *testing.T) {
	config.ReadConfig("")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		region     string
		actID      string
		size       int
		startIndex int
	}
	var tests []struct {
		name    string
		args    args
		want    *LeaderboardDto
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetLeaderboard(tt.args.region, tt.args.actID, tt.args.size, tt.args.startIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLeaderboard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLeaderboard() got = %v, want %v", got, tt.want)
			}
		})
	}
}
