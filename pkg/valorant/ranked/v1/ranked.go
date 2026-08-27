package v1

import (
	"encoding/json"
	"net/url"
	"strconv"

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

func (c *Client) GetLeaderboard(region string, actID string, size int, startIndex int) (*LeaderboardDto, error) {
	path := "/val/ranked/v1/leaderboards/by-act/" + url.PathEscape(actID) + c.c.MakeQueryString(
		map[string]string{"size": strconv.Itoa(size), "startIndex": strconv.Itoa(startIndex)})
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var l LeaderboardDto
	err = json.Unmarshal(res, &l)
	return &l, err
}
