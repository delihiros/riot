package v3

import (
	"encoding/json"

	"github.com/delihiros/riot/pkg/client"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) GetChampionInfo(region string) (*ChampionInfo, error) {
	res, err := c.c.SimpleGet(region, "/lol/platform/v3/champion-rotations")
	if err != nil {
		return nil, err
	}
	var info ChampionInfo
	if err := json.Unmarshal(res, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
