package camera

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Credential struct {
	// Host is the camera address, optionally with a scheme and port:
	// "192.168.1.64", "192.168.1.64:8080" (HTTP) or "https://192.168.1.64[:8443]".
	// Without a scheme, HTTP is used. The ONVIF driver always uses HTTP; see
	// the onvif package for how it interprets an https:// host.
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

	// TLSConfig configures HTTPS connections to cameras whose Host uses the
	// https:// scheme. Cameras usually ship self-signed certificates, so set
	// RootCAs to trust the camera's certificate (recommended), or
	// InsecureSkipVerify to skip verification (vulnerable to MITM).
	// Nil verifies certificates against the system roots.
	TLSConfig *tls.Config
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
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if cfg.TLSConfig != nil {
		transport.TLSClientConfig = cfg.TLSConfig.Clone()
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
