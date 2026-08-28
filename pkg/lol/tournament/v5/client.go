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

func New(c *client.Client) *Client {
	return &Client{c: c}
}

func (c *Client) CreateTournamentCode(region string, tournamentID int64, count *int, params TournamentCodeParametersV5) ([]string, error) {
	query := url.Values{"tournamentId": {strconv.FormatInt(tournamentID, 10)}}
	if count != nil {
		query.Set("count", strconv.Itoa(*count))
	}
	var codes []string
	if err := c.post(region, "/lol/tournament/v5/codes?"+query.Encode(), params, &codes); err != nil {
		return nil, err
	}
	return codes, nil
}

func (c *Client) GetTournamentCode(region, tournamentCode string) (*TournamentCodeV5DTO, error) {
	var code TournamentCodeV5DTO
	if err := c.get(region, "/lol/tournament/v5/codes/"+url.PathEscape(tournamentCode), &code); err != nil {
		return nil, err
	}
	return &code, nil
}

func (c *Client) UpdateCode(region, tournamentCode string, params *TournamentCodeUpdateParametersV5) error {
	var body []byte
	var err error
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return err
		}
	}
	_, err = c.c.SimplePut(region, "/lol/tournament/v5/codes/"+url.PathEscape(tournamentCode), body)
	return err
}

func (c *Client) GetGames(region, tournamentCode string) ([]TournamentGamesV5, error) {
	var games []TournamentGamesV5
	if err := c.get(region, "/lol/tournament/v5/games/by-code/"+url.PathEscape(tournamentCode), &games); err != nil {
		return nil, err
	}
	return games, nil
}

func (c *Client) GetLobbyEventsByCode(region, tournamentCode string) (*LobbyEventV5DTOWrapper, error) {
	var events LobbyEventV5DTOWrapper
	if err := c.get(region, "/lol/tournament/v5/lobby-events/by-code/"+url.PathEscape(tournamentCode), &events); err != nil {
		return nil, err
	}
	return &events, nil
}

func (c *Client) RegisterProviderData(region string, params ProviderRegistrationParametersV5) (int, error) {
	var providerID int
	if err := c.post(region, "/lol/tournament/v5/providers", params, &providerID); err != nil {
		return 0, err
	}
	return providerID, nil
}

func (c *Client) RegisterTournament(region string, params TournamentRegistrationParametersV5) (int, error) {
	var tournamentID int
	if err := c.post(region, "/lol/tournament/v5/tournaments", params, &tournamentID); err != nil {
		return 0, err
	}
	return tournamentID, nil
}

func (c *Client) get(region, path string, dst any) error {
	res, err := c.c.SimpleGet(region, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}

func (c *Client) post(region, path string, params, dst any) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	res, err := c.c.SimplePost(region, path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dst)
}
