package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	apiKey string
	region string
	client *http.Client
}

func New(apiKey string) *Client {
	client := &http.Client{}
	return &Client{
		apiKey: apiKey,
		client: client,
	}
}

func (c *Client) SimpleGet(region string, url string) ([]byte, error) {
	params := map[string]string{
		"X-Riot-Token": c.apiKey,
		"Content-Type": "application/json;charset=UTF-8",
	}
	requestingURL := "https://" + region + ".api.riotgames.com" + url
	return c.GetWithHeaderParams(requestingURL, params)
}

func (c *Client) GetWithHeaderParams(url string, params map[string]string) ([]byte, error) {
	return c.GetWithHeadersAndQueries(url, params, map[string]string{})
}

func (c *Client) GetWithURLQueries(url string, queries map[string]string) ([]byte, error) {
	return c.GetWithHeadersAndQueries(url, map[string]string{}, queries)
}

func (c *Client) GetWithHeadersAndQueries(url string, headerParams map[string]string, queries map[string]string) ([]byte, error) {
	queryString := c.MakeQueryString(queries)
	req, err := http.NewRequest("GET", url+queryString, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headerParams {
		req.Header.Set(k, v)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) MakeQueryString(params map[string]string) string {
	u := &url.URL{}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
