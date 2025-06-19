package sbdb

import (
	"fmt"
	"net/http"
	"net/url"
)

const Endpoint = "https://ssd-api.jpl.nasa.gov/sbdb_query.api"

type Client struct {
	http.Client
	Endpoint string
}

func DefaultClient() *Client {
	return &Client{
		Endpoint: Endpoint,
	}
}

func (c *Client) Get(f Filter) (*http.Response, error) {
	u, err := c.GetURL(f)
	if err != nil {
		return nil, err
	}

	return c.Client.Get(u.String())
}

func (c *Client) GetURL(f Filter) (*url.URL, error) {
	v, err := f.Values()
	if err != nil {
		return nil, fmt.Errorf("error parsing filter: %w", err)
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing endpoint: %w", err)
	}
	u.RawQuery = v.Encode()
	return u, nil
}
