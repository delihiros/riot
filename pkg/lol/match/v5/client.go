package v5

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/delihiros/riot/pkg/client"
)

type Client struct {
	c *client.Client
}

type MatchIDsOptions struct {
	StartTime *int64
	EndTime   *int64
	Queue     *int
	Type      *string
	Start     *int
	Count     *int
}

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) GetMatchIDsByPUUID(region, puuid string, options *MatchIDsOptions) ([]string, error) {
	path := "/lol/match/v5/matches/by-puuid/" + url.PathEscape(puuid) + "/ids"
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
	if err := c.get(region, path, &matchIDs); err != nil {
		return nil, err
	}
	return matchIDs, nil
}

func (c *Client) GetReplay(region, puuid string) (*ReplayDTO, error) {
	var replay ReplayDTO
	if err := c.get(region, "/lol/match/v5/matches/by-puuid/"+url.PathEscape(puuid)+"/replays", &replay); err != nil {
		return nil, err
	}
	return &replay, nil
}

func (c *Client) GetMatch(region, matchID string) (*MatchDto, error) {
	var match MatchDto
	if err := c.get(region, "/lol/match/v5/matches/"+url.PathEscape(matchID), &match); err != nil {
		return nil, err
	}
	return &match, nil
}

func (c *Client) GetTimeline(region, matchID string) (*TimelineDto, error) {
	var timeline TimelineDto
	if err := c.get(region, "/lol/match/v5/matches/"+url.PathEscape(matchID)+"/timeline", &timeline); err != nil {
		return nil, err
	}
	return &timeline, nil
}

func (c *Client) get(region, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
