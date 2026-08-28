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

func (c *Client) GetAllChampionMasteriesByPUUID(region string, puuid string) ([]ChampionMasteryDto, error) {
	var masteries []ChampionMasteryDto
	err := c.get(region, "/lol/champion-mastery/v4/champion-masteries/by-puuid/"+url.PathEscape(puuid), &masteries)
	return masteries, err
}

func (c *Client) GetChampionMasteryByPUUID(region string, puuid string, championID int64) (*ChampionMasteryDto, error) {
	var mastery ChampionMasteryDto
	err := c.get(region, "/lol/champion-mastery/v4/champion-masteries/by-puuid/"+url.PathEscape(puuid)+"/by-champion/"+strconv.FormatInt(championID, 10), &mastery)
	return &mastery, err
}

func (c *Client) GetTopChampionMasteriesByPUUID(region string, puuid string, count *int) ([]ChampionMasteryDto, error) {
	path := "/lol/champion-mastery/v4/champion-masteries/by-puuid/" + url.PathEscape(puuid) + "/top"
	if count != nil {
		path += "?count=" + strconv.Itoa(*count)
	}
	var masteries []ChampionMasteryDto
	err := c.get(region, path, &masteries)
	return masteries, err
}

func (c *Client) GetChampionMasteryScoreByPUUID(region string, puuid string) (int, error) {
	var score int
	err := c.get(region, "/lol/champion-mastery/v4/scores/by-puuid/"+url.PathEscape(puuid), &score)
	return score, err
}

func (c *Client) get(region string, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
