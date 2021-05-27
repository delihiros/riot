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

func (c *Client) Get(region string, url string) ([]byte, error) {
	params := map[string]string{}
	requestingURL := "https://" + region + ".api.riotgames.com" + url
	params["X-Riot-Token"] = c.apiKey
	params["Content-Type"] = "application/json;charset=UTF-8"
	return c.getWithHeaderParams(requestingURL, params)
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

func (c *Client) getWithHeaderParams(url string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range params {
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
