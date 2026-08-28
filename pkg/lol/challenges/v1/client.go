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

func (c *Client) GetAllChallengeConfigs(region string) ([]ChallengeConfigInfoDto, error) {
	var configs []ChallengeConfigInfoDto
	if err := c.get(region, "/lol/challenges/v1/challenges/config", &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (c *Client) GetAllChallengePercentiles(region string) (map[int64]map[int]map[string]float64, error) {
	var percentiles map[int64]map[int]map[string]float64
	if err := c.get(region, "/lol/challenges/v1/challenges/percentiles", &percentiles); err != nil {
		return nil, err
	}
	return percentiles, nil
}

func (c *Client) GetChallengeConfigs(region string, challengeID int64) (*ChallengeConfigInfoDto, error) {
	var config ChallengeConfigInfoDto
	path := "/lol/challenges/v1/challenges/" + strconv.FormatInt(challengeID, 10) + "/config"
	if err := c.get(region, path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Client) GetChallengeLeaderboards(region string, challengeID int64, level string, limit *int) ([]ApexPlayerInfoDto, error) {
	path := "/lol/challenges/v1/challenges/" + strconv.FormatInt(challengeID, 10) + "/leaderboards/by-level/" + url.PathEscape(level)
	if limit != nil {
		path += "?" + url.Values{"limit": {strconv.Itoa(*limit)}}.Encode()
	}
	var players []ApexPlayerInfoDto
	if err := c.get(region, path, &players); err != nil {
		return nil, err
	}
	return players, nil
}

func (c *Client) GetChallengePercentiles(region string, challengeID int64) (map[string]float64, error) {
	var percentiles map[string]float64
	path := "/lol/challenges/v1/challenges/" + strconv.FormatInt(challengeID, 10) + "/percentiles"
	if err := c.get(region, path, &percentiles); err != nil {
		return nil, err
	}
	return percentiles, nil
}

func (c *Client) GetPlayerData(region, puuid string) (*PlayerInfoDto, error) {
	var player PlayerInfoDto
	if err := c.get(region, "/lol/challenges/v1/player-data/"+url.PathEscape(puuid), &player); err != nil {
		return nil, err
	}
	return &player, nil
}

func (c *Client) get(region, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
