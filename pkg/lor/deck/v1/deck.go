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

func (c *Client) GetUserDeck(region string, authorization string) ([]*DeckDto, error) {
	res, err := c.c.GetWithHeaderParams("/lor/deck/v1/decks/me", map[string]string{"Authorization": authorization})
	if err != nil {
		return nil, err
	}
	var dd []*DeckDto
	err = json.Unmarshal(res, dd)
	return dd, err
}
