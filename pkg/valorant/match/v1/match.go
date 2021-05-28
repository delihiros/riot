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

func (c *Client) GetMatchByID(region string, matchID string) (*MatchDto, error) {
	url := "/val/match/v1/matches/" + matchID
	res, err := c.c.SimpleGet(region, url)
	if err != nil {
		return nil, err
	}
	var md MatchDto
	err = json.Unmarshal(res, &md)
	return &md, err
}

func (c *Client) GetMatchListByPUUID(region string, puuid string) (*MatchListDto, error) {
	url := "/val/match/v1/matchlists/by-puuid/" + puuid
	res, err := c.c.SimpleGet(region, url)
	if err != nil {
		return nil, err
	}
	var md MatchListDto
	err = json.Unmarshal(res, &md)
	return &md, err
}
