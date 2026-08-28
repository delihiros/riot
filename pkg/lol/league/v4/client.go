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

func (c *Client) GetChallengerLeague(region, queue string) (*LeagueListDTO, error) {
	return c.getLeague(region, "/lol/league/v4/challengerleagues/by-queue/"+url.PathEscape(queue))
}

func (c *Client) GetLeagueEntriesByPUUID(region, puuid string) ([]LeagueEntryDTO, error) {
	var entries []LeagueEntryDTO
	if err := c.get(region, "/lol/league/v4/entries/by-puuid/"+url.PathEscape(puuid), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) GetLeagueEntries(region, queue, tier, division string, page *int) ([]LeagueEntryDTO, error) {
	path := "/lol/league/v4/entries/" + url.PathEscape(queue) + "/" + url.PathEscape(tier) + "/" + url.PathEscape(division)
	if page != nil {
		path += "?page=" + strconv.Itoa(*page)
	}
	var entries []LeagueEntryDTO
	if err := c.get(region, path, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) GetGrandmasterLeague(region, queue string) (*LeagueListDTO, error) {
	return c.getLeague(region, "/lol/league/v4/grandmasterleagues/by-queue/"+url.PathEscape(queue))
}

func (c *Client) GetMasterLeague(region, queue string) (*LeagueListDTO, error) {
	return c.getLeague(region, "/lol/league/v4/masterleagues/by-queue/"+url.PathEscape(queue))
}

func (c *Client) getLeague(region, path string) (*LeagueListDTO, error) {
	var league LeagueListDTO
	if err := c.get(region, path, &league); err != nil {
		return nil, err
	}
	return &league, nil
}

func (c *Client) get(region, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
