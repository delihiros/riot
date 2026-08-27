package v1

import (
	"encoding/json"

	"github.com/delihiros/riot/pkg/client"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) GetLeaderboards(region string) (*LeaderboardDto, error) {
	res, err := c.c.SimpleGet(region, "/lor/ranked/v1/leaderboards")
	if err != nil {
		return nil, err
	}
	var ld LeaderboardDto
	err = json.Unmarshal(res, &ld)
	return &ld, err
}
