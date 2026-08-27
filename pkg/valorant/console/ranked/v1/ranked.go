package v1

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/delihiros/riot/pkg/client"
	pcranked "github.com/delihiros/riot/pkg/valorant/ranked/v1"
)

type Client struct {
	c *client.Client
}

type PlayerDto = pcranked.PlayerDto

type LeaderboardDto struct {
	ActID        string            `json:"actId"`
	TotalPlayers int64             `json:"totalPlayers"`
	Query        string            `json:"query"`
	Shard        string            `json:"shard"`
	Players      []*PlayerDto      `json:"players"`
	TierDetails  []json.RawMessage `json:"tierDetails"`
}

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) GetLeaderboard(region string, actID string, platformType string, size int, startIndex int) (*LeaderboardDto, error) {
	path := "/val/console/ranked/v1/leaderboards/by-act/" + url.PathEscape(actID)
	path += c.c.MakeQueryString(map[string]string{
		"platformType": platformType,
		"size":         strconv.Itoa(size),
		"startIndex":   strconv.Itoa(startIndex),
	})
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var leaderboard LeaderboardDto
	err = json.Unmarshal(res, &leaderboard)
	return &leaderboard, err
}
