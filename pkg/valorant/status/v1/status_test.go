package v1

import (
	"reflect"
	"riot/config"
	"riot/pkg/client"
	"testing"
)

func TestClient_GetStatus(t *testing.T) {
	config.ReadConfig("")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		region string
	}
	var tests []struct {
		name    string
		args    args
		want    *PlatformDataDto
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetStatus(tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
