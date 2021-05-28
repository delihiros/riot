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

func (c *Client) GetStatus(region string) (*PlatformDataDto, error) {
	res, err := c.c.SimpleGet(region, "/val/status/v1/platform-data")
	if err != nil {
		return nil, err
	}
	var pd PlatformDataDto
	err = json.Unmarshal(res, &pd)
	return &pd, err
}
