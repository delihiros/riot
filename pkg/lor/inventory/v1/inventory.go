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

func (c *Client) GetInventory(region string, authorization string) ([]*CardDto, error) {
	res, err := c.c.GetWithHeaderParams("/lor/inventory/v1/cards/me", map[string]string{"Authorization": authorization, "region": region})
	if err != nil {
		return nil, err
	}
	var cd []*CardDto
	err = json.Unmarshal(res, cd)
	return cd, err
}
