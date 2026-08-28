package v4

import (
	"encoding/json"

	"github.com/delihiros/riot/pkg/client"
)

type Client struct {
	c *client.Client
}

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) GetPlatformData(region string) (*PlatformDataDto, error) {
	res, err := c.c.SimpleGet(region, "/lol/status/v4/platform-data")
	if err != nil {
		return nil, err
	}
	var data PlatformDataDto
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
