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
	res, err := c.c.GetWithHeaderParams("/lor/status/v1/platform-data", map[string]string{"region": region})
	if err != nil {
		return nil, err
	}
	var pdd PlatformDataDto
	err = json.Unmarshal(res, &pdd)
	return &pdd, err
}
