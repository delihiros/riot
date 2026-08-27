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

func (c *Client) GetContent(region string, locale string) (*ContentDto, error) {
	path := "/val/content/v1/contents"
	if locale != "" {
		path += c.c.MakeQueryString(map[string]string{"locale": locale})
	}
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return nil, err
	}
	var cd ContentDto
	err = json.Unmarshal(res, &cd)
	return &cd, err
}
