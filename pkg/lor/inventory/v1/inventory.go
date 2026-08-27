package v1

import (
	"encoding/json"

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

func (c *Client) GetInventory(region string, authorization string) ([]*CardDto, error) {
	res, err := c.c.GetWithRegionAndHeaders(region, "/lor/inventory/v1/cards/me", map[string]string{"Authorization": authorization})
	if err != nil {
		return nil, err
	}
	var cd []*CardDto
	err = json.Unmarshal(res, &cd)
	return cd, err
}
