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
	return &Client{c: c}
}

func (c *Client) GetPlayersByPUUID(region, puuid string) ([]PlayerDto, error) {
	var players []PlayerDto
	if err := c.get(region, "/lol/clash/v1/players/by-puuid/"+url.PathEscape(puuid), &players); err != nil {
		return nil, err
	}
	return players, nil
}

func (c *Client) GetTeamByID(region, teamID string) (*TeamDto, error) {
	var team TeamDto
	if err := c.get(region, "/lol/clash/v1/teams/"+url.PathEscape(teamID), &team); err != nil {
		return nil, err
	}
	return &team, nil
}

func (c *Client) GetTournaments(region string) ([]TournamentDto, error) {
	var tournaments []TournamentDto
	if err := c.get(region, "/lol/clash/v1/tournaments", &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (c *Client) GetTournamentByTeam(region, teamID string) (*TournamentDto, error) {
	var tournament TournamentDto
	if err := c.get(region, "/lol/clash/v1/tournaments/by-team/"+url.PathEscape(teamID), &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (c *Client) GetTournamentByID(region string, tournamentID int) (*TournamentDto, error) {
	var tournament TournamentDto
	if err := c.get(region, "/lol/clash/v1/tournaments/"+strconv.Itoa(tournamentID), &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (c *Client) get(region, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
