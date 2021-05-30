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

func (c *Client) GetMatchListByPUUID(region string, puuid string) ([]string, error) {
	url := "/lor/match/v1/matches/by-puuid/" + puuid + "/ids"
	res, err := c.c.GetWithHeaderParams(url, map[string]string{"region": region})
	if err != nil {
		return nil, err
	}
	var matchIDs []string
	err = json.Unmarshal(res, matchIDs) // TODO: need verification
	return matchIDs, err
}

func (c *Client) GetMatchByID(region string, matchID string) (*MatchDto, error) {
	url := "/lor/match/v1/matches/" + matchID
	res, err := c.c.GetWithHeaderParams(url, map[string]string{"region": region})
	if err != nil {
		return nil, err
	}
	var m MatchDto
	err = json.Unmarshal(res, &m)
	return &m, err
}
