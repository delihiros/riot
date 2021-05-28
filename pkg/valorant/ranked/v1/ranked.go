package v1

import (
	"encoding/json"
	"riot/pkg/client"
	"strconv"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) GetLeaderboard(region string, actID string, size int, startIndex int) (*LeaderboardDto, error) {
	url := "/val/ranked/v1/leaderboards/by-act/" + actID + c.c.MakeQueryString(
		map[string]string{"size": strconv.Itoa(size), "startIndex": strconv.Itoa(startIndex)})
	res, err := c.c.SimpleGet(region, url)
	if err != nil {
		return nil, err
	}
	var l LeaderboardDto
	err = json.Unmarshal(res, &l)
	return &l, err
}
