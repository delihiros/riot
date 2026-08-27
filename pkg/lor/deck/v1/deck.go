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

func (c *Client) GetUserDeck(region string, authorization string) ([]*DeckDto, error) {
	res, err := c.c.GetWithRegionAndHeaders(region, "/lor/deck/v1/decks/me", map[string]string{"Authorization": authorization})
	if err != nil {
		return nil, err
	}
	var dd []*DeckDto
	err = json.Unmarshal(res, &dd)
	return dd, err
}

func (c *Client) CreateDeck(region string, authorization string, deck *NewDeckDto) (string, error) {
	body, err := json.Marshal(deck)
	if err != nil {
		return "", err
	}
	res, err := c.c.PostWithRegionAndHeaders(region, "/lor/deck/v1/decks/me", map[string]string{
		"Authorization": authorization,
		"Content-Type":  "application/json",
	}, body)
	if err != nil {
		return "", err
	}
	var id string
	err = json.Unmarshal(res, &id)
	return id, err
}
