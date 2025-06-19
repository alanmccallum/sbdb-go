package sbdb

import (
	"fmt"
	"net/http"
	"net/url"
)

const Endpoint = "https://ssd-api.jpl.nasa.gov/sbdb_query.api"

// Client wraps http.Client
// The Client.Endpoint field can be set to use a custom api server.
type Client struct {
	http.Client
	Endpoint string
}

// Get issues a GET to the Client.Endpoint using the specified Filter
func (c *Client) Get(f Filter) (*http.Response, error) {
	u, err := c.GetURL(f)
	if err != nil {
		return nil, err
	}

	return c.Client.Get(u.String())
}

// GetURL builds a *url.URL from the Client.Endpoint and the specified Filter
func (c *Client) GetURL(f Filter) (*url.URL, error) {
	v, err := f.Values()
	if err != nil {
		return nil, fmt.Errorf("error parsing filter: %w", err)
	}
	ep := Endpoint
	if c.Endpoint != "" {
		ep = c.Endpoint
	}
	u, err := url.Parse(ep)
	if err != nil {
		return nil, fmt.Errorf("error parsing endpoint: %w", err)
	}
	u.RawQuery = v.Encode()
	return u, nil
}
