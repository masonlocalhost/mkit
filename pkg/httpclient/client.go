package httpclient

import (
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	config *Config
}

func NewClient(cfg *Config) *Client {
	if cfg.DefaultTimeout == 0 {
		cfg.DefaultTimeout = 30 * time.Second
	}

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				// customize transport like connection settings, tls...
			},
		},
		config: cfg,
	}
}

func (c *Client) Req() *Request {
	return &Request{
		headers:      make(map[string][]string),
		queries:      make(map[string]string),
		client:       c,
		timeout:      c.config.DefaultTimeout,
		retryCount:   c.config.DefaultRetryCount,
		requestBody:  nil,
		responseBody: nil,
	}
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.client
}
