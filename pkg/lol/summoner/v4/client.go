package v4

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

func (c *Client) GetByPUUID(region, puuid string) (*SummonerDTO, error) {
	res, err := c.c.SimpleGet(region, "/lol/summoner/v4/summoners/by-puuid/"+url.PathEscape(puuid))
	if err != nil {
		return nil, err
	}
	var summoner SummonerDTO
	if err := json.Unmarshal(res, &summoner); err != nil {
		return nil, err
	}
	return &summoner, nil
}

func (c *Client) GetByAccessToken(region, accessToken string) (*SummonerDTO, error) {
	res, err := c.c.GetWithRegionAndHeaders(region, "/lol/summoner/v4/summoners/me", map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return nil, err
	}
	var summoner SummonerDTO
	if err := json.Unmarshal(res, &summoner); err != nil {
		return nil, err
	}
	return &summoner, nil
}
