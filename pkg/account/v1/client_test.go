package v1

import (
	"reflect"
	"riot/config"
	"riot/pkg/client"
	"testing"
)

func TestClient_GetAccountByPUUID(t *testing.T) {
	config.ReadConfig("../../../config/test.yaml")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		puuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *AccountDto
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				puuid: cfg.GetString("test.puuid"),
			},
			want: &AccountDto{
				PUUID:    cfg.GetString("test.puuid"),
				GameName: cfg.GetString("test.game_name"),
				TagLine:  cfg.GetString("test.tag_line"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetAccountByPUUID("asia", tt.args.puuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByPUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountByPUUID2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccountByRiotID(t *testing.T) {
	config.ReadConfig("../../../config/test.yaml")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		puuid    string
		gameName string
		tagLine  string
	}
	tests := []struct {
		name    string
		args    args
		want    *AccountDto
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				gameName: cfg.GetString("test.game_name"),
				tagLine:  cfg.GetString("test.tag_line"),
			},
			want: &AccountDto{
				PUUID:    cfg.GetString("test.puuid"),
				GameName: cfg.GetString("test.game_name"),
				TagLine:  cfg.GetString("test.tag_line"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetAccountByRiotID("asia", tt.args.gameName, tt.args.tagLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByRiotID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountByRiotID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetActiveShard(t *testing.T) {
	config.ReadConfig("../../../config/test.yaml")
	cfg := config.GetConfig()
	http := client.New(cfg.GetString("riot.api_key"))

	type args struct {
		game  string
		puuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *ActiveShardDto
		wantErr bool
	}{
		{
			name: "val-ap",
			args: args{
				game:  "val",
				puuid: cfg.GetString("test.puuid"),
			},
			want: &ActiveShardDto{
				PUUID:       cfg.GetString("test.puuid"),
				Game:        "val",
				ActiveShard: "ap",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(http)
			got, err := c.GetActiveShard("asia", tt.args.game, tt.args.puuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveShard() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveShard() got = %v, want %v", got, tt.want)
			}
		})
	}
}
