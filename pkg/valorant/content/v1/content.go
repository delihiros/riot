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

func (c *Client) GetContent(region string, locale string) (*ContentDto, error) {
	url := "/val/content/v1/contents" + c.c.MakeQueryString(map[string]string{"locale": locale})
	res, err := c.c.Get(region, url)
	if err != nil {
		return nil, err
	}
	var cd ContentDto
	err = json.Unmarshal(res, &cd)
	return &cd, err
}
