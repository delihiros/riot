package v1

import (
	"encoding/json"
	"net/url"

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

func (c *Client) GetAccountByPUUID(region string, puuid string) (*AccountDto, error) {
	res, err := c.c.SimpleGet(region, "/riot/account/v1/accounts/by-puuid/"+url.PathEscape(puuid))
	if err != nil {
		return nil, err
	}
	var a AccountDto
	err = json.Unmarshal(res, &a)
	return &a, err
}

func (c *Client) GetAccountByRiotID(region string, gameName string, tagLine string) (*AccountDto, error) {
	res, err := c.c.SimpleGet(region, "/riot/account/v1/accounts/by-riot-id/"+url.PathEscape(gameName)+"/"+url.PathEscape(tagLine))
	if err != nil {
		return nil, err
	}
	var a AccountDto
	err = json.Unmarshal(res, &a)
	return &a, err
}

func (c *Client) GetActiveShard(region string, game string, puuid string) (*ActiveShardDto, error) {
	res, err := c.c.SimpleGet(region, "/riot/account/v1/active-shards/by-game/"+url.PathEscape(game)+"/by-puuid/"+url.PathEscape(puuid))
	if err != nil {
		return nil, err
	}
	var as ActiveShardDto
	err = json.Unmarshal(res, &as)
	return &as, err
}

func (c *Client) GetActiveRegion(region string, game string, puuid string) (*AccountRegionDto, error) {
	res, err := c.c.SimpleGet(region, "/riot/account/v1/region/by-game/"+url.PathEscape(game)+"/by-puuid/"+url.PathEscape(puuid))
	if err != nil {
		return nil, err
	}
	var accountRegion AccountRegionDto
	err = json.Unmarshal(res, &accountRegion)
	return &accountRegion, err
}

func (c *Client) GetAccountByAccessToken(region string, accessToken string) (*AccountDto, error) {
	res, err := c.c.GetWithRegionAndHeaders(region, "/riot/account/v1/accounts/me", map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return nil, err
	}
	var account AccountDto
	err = json.Unmarshal(res, &account)
	return &account, err
}
