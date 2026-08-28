package v1

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

func (c *Client) GetMatchIDs(region, accessToken string, options *MatchIDsOptions) ([]string, error) {
	path := "/lol/rso-match/v1/matches/ids"
	if options != nil {
		query := url.Values{}
		if options.StartTime != nil {
			query.Set("startTime", strconv.FormatInt(*options.StartTime, 10))
		}
		if options.EndTime != nil {
			query.Set("endTime", strconv.FormatInt(*options.EndTime, 10))
		}
		if options.Queue != nil {
			query.Set("queue", strconv.Itoa(*options.Queue))
		}
		if options.Type != nil {
			query.Set("type", *options.Type)
		}
		if options.Start != nil {
			query.Set("start", strconv.Itoa(*options.Start))
		}
		if options.Count != nil {
			query.Set("count", strconv.Itoa(*options.Count))
		}
		if len(query) != 0 {
			path += "?" + query.Encode()
		}
	}

	var matchIDs []string
	if err := c.get(region, accessToken, path, &matchIDs); err != nil {
		return nil, err
	}
	return matchIDs, nil
}

func (c *Client) GetMatch(region, accessToken, matchID string) (*MatchDto, error) {
	var match MatchDto
	if err := c.get(region, accessToken, "/lol/rso-match/v1/matches/"+url.PathEscape(matchID), &match); err != nil {
		return nil, err
	}
	return &match, nil
}

func (c *Client) GetTimeline(region, accessToken, matchID string) (*TimelineDto, error) {
	var timeline TimelineDto
	if err := c.get(region, accessToken, "/lol/rso-match/v1/matches/"+url.PathEscape(matchID)+"/timeline", &timeline); err != nil {
		return nil, err
	}
	return &timeline, nil
}

func (c *Client) get(region, accessToken, path string, dst any) error {
	res, err := c.c.GetWithRegionAndHeaders(region, path, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
