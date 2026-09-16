package camera

import (
	"net/http"
	"time"
)

type Credential struct {
	Host     string
	Username string
	Password string
}

type Config struct {
	Timeout         time.Duration
	DownloadTimeout time.Duration
}

type Client struct {
	http         *http.Client
	downloadHTTP *http.Client
}

func New(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.DownloadTimeout == 0 {
		cfg.DownloadTimeout = 5 * time.Minute
	}

	transport := &http.Transport{
		DisableKeepAlives: false,

		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	return &Client{
		http: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		downloadHTTP: &http.Client{
			Timeout:   cfg.DownloadTimeout,
			Transport: transport,
		},
	}
}

func (c *Client) HTTP() *http.Client {
	return c.http
}

func (c *Client) HTTPDownload() *http.Client {
	return c.downloadHTTP
}
