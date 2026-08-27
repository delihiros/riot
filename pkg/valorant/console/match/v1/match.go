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
	return &Client{c: c}
}

func (c *Client) GetMatchByID(region string, matchID string) (*MatchDto, error) {
	path := "/val/match/console/v1/matches/" + url.PathEscape(matchID)
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var match MatchDto
	err = json.Unmarshal(res, &match)
	return &match, err
}

func (c *Client) GetMatchListByPUUID(region string, puuid string, platformType string) (*MatchListDto, error) {
	path := "/val/match/console/v1/matchlists/by-puuid/" + url.PathEscape(puuid)
	path += c.c.MakeQueryString(map[string]string{"platformType": platformType})
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var matches MatchListDto
	err = json.Unmarshal(res, &matches)
	return &matches, err
}

func (c *Client) GetRecentMatches(region string, queue string) (*RecentMatchesDto, error) {
	path := "/val/match/console/v1/recent-matches/by-queue/" + url.PathEscape(queue)
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var matches RecentMatchesDto
	err = json.Unmarshal(res, &matches)
	return &matches, err
}
