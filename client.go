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
	// Timeout is the HTTP timeout for regular requests. Default: 30s.
	Timeout time.Duration

	// DownloadTimeout is the HTTP timeout for recording downloads. Default: 5m.
	DownloadTimeout time.Duration

	// MoveTimeout is how long PTZ moves that wait for the camera to reach a
	// target (e.g. SetAbsolute, SetZoom) keep polling before giving up.
	// A context deadline, if shorter, still applies. Default: 5s.
	MoveTimeout time.Duration
}

type Client struct {
	http         *http.Client
	downloadHTTP *http.Client
	moveTimeout  time.Duration
}

func New(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.DownloadTimeout == 0 {
		cfg.DownloadTimeout = 5 * time.Minute
	}
	if cfg.MoveTimeout == 0 {
		cfg.MoveTimeout = 5 * time.Second
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
		moveTimeout: cfg.MoveTimeout,
	}
}

func (c *Client) HTTP() *http.Client {
	return c.http
}

func (c *Client) HTTPDownload() *http.Client {
	return c.downloadHTTP
}

// MoveTimeout returns how long PTZ moves wait for the camera to reach a target.
func (c *Client) MoveTimeout() time.Duration {
	return c.moveTimeout
}
