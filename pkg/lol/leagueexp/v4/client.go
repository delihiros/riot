package v4

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
	return &Client{c: c}
}

func (c *Client) GetLeagueEntries(region, queue, tier, division string, page *int) ([]LeagueEntryDTO, error) {
	path := "/lol/league-exp/v4/entries/" + url.PathEscape(queue) + "/" + url.PathEscape(tier) + "/" + url.PathEscape(division)
	if page != nil {
		query := url.Values{"page": {strconv.Itoa(*page)}}
		path += "?" + query.Encode()
	}

	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var entries []LeagueEntryDTO
	if err := json.Unmarshal(res, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
