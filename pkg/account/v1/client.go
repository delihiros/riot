package v1

import (
	"encoding/json"
	"riot/pkg/client"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) GetAccountByPUUID(region string, puuid string) (*AccountDto, error) {
	res, err := c.c.Get(region, "/riot/account/v1/accounts/by-puuid/"+puuid)
	if err != nil {
		return nil, err
	}
	var a AccountDto
	err = json.Unmarshal(res, &a)
	return &a, err
}

func (c *Client) GetAccountByRiotID(region string, gameName string, tagLine string) (*AccountDto, error) {
	res, err := c.c.Get(region, "/riot/account/v1/accounts/by-riot-id/"+gameName+"/"+tagLine)
	if err != nil {
		return nil, err
	}
	var a AccountDto
	err = json.Unmarshal(res, &a)
	return &a, err
}

func (c *Client) GetActiveShard(region string, game string, puuid string) (*ActiveShardDto, error) {
	res, err := c.c.Get(region, "/riot/account/v1/active-shards/by-game/"+game+"/by-puuid/"+puuid)
	if err != nil {
		return nil, err
	}
	var as ActiveShardDto
	err = json.Unmarshal(res, &as)
	return &as, err
}
