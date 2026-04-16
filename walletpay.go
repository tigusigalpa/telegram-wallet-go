package walletpay

import (
	"net/http"
	"time"
)

// Client is the Wallet Pay API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Wallet Pay client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: "https://pay.wallet.tg",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}
