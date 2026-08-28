package v5

import (
	"encoding/json"
	"net/url"

	"github.com/delihiros/riot/pkg/client"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) GetCurrentGameInfoByPUUID(region, puuid string) (*CurrentGameInfo, error) {
	res, err := c.c.SimpleGet(region, "/lol/spectator/v5/active-games/by-summoner/"+url.PathEscape(puuid))
	if err != nil {
		return nil, err
	}
	var game CurrentGameInfo
	if err := json.Unmarshal(res, &game); err != nil {
		return nil, err
	}
	return &game, nil
}
