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

func (c *Client) GetMatchListByPUUID(region string, puuid string) ([]string, error) {
	path := "/lor/match/v1/matches/by-puuid/" + url.PathEscape(puuid) + "/ids"
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var matchIDs []string
	err = json.Unmarshal(res, &matchIDs)
	return matchIDs, err
}

func (c *Client) GetMatchByID(region string, matchID string) (*MatchDto, error) {
	path := "/lor/match/v1/matches/" + url.PathEscape(matchID)
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var m MatchDto
	err = json.Unmarshal(res, &m)
	return &m, err
}
