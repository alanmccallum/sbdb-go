package sbdb

import (
	"fmt"
	"net/http"
	"net/url"
)

// Endpoint is the default base URL for the SBDB Query API.
// It can be overridden via Client.Endpoint for testing or custom servers.
const Endpoint = "https://ssd-api.jpl.nasa.gov/sbdb_query.api"

// Client wraps http.Client and provides helpers for interacting with
// the SBDB Query API. The Client.Endpoint field can be set to use a
// custom API server.
type Client struct {
	http.Client
	Endpoint string
}

// Get issues a GET request using the provided Filter.
// The request is sent to Endpoint or Client.Endpoint if set.
func (c *Client) Get(f Filter) (*http.Response, error) {
	u, err := c.GetURL(f)
	if err != nil {
		return nil, err
	}

	return c.Client.Get(u.String())
}

// GetURL builds a URL for the request represented by the Filter. If
// Client.Endpoint is empty, the default Endpoint constant is used.
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
