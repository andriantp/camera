package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/icholy/digest"
)

type credentialKey struct {
	username string
	password string
}

// authClients holds the digest-authenticated clients for one credential.
type authClients struct {
	http     *http.Client
	download *http.Client
}

// authClientsFor returns cached digest clients for the camera's credential.
// digest.Transport caches the server challenge per host, so reusing it lets
// subsequent requests authenticate in one round trip instead of two
// (unauthenticated request -> 401 -> retry).
func (c *Client) authClientsFor(camera Camera) *authClients {
	key := credentialKey{username: camera.Username, password: camera.Password}

	if v, ok := c.auth.Load(key); ok {
		return v.(*authClients)
	}

	newClient := func(base *http.Client) *http.Client {
		return &http.Client{
			Timeout: base.Timeout,
			Transport: &digest.Transport{
				Username:  camera.Username,
				Password:  camera.Password,
				Transport: base.Transport,
			},
		}
	}

	v, _ := c.auth.LoadOrStore(key, &authClients{
		http:     newClient(c.http),
		download: newClient(c.downloadHTTP),
	})

	return v.(*authClients)
}

func (c *Client) request(ctx context.Context, method string, camera Camera, path string, body io.Reader) (*http.Response, error) {
	resp, err := c.do(ctx, c.authClientsFor(camera).http, method, camera, path, body)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) requestDownload(ctx context.Context, method string, camera Camera, path string, body io.Reader) (*http.Response, error) {
	resp, err := c.do(ctx, c.authClientsFor(camera).download, method, camera, path, body)
	if err != nil {
		return nil, fmt.Errorf("do download request: %w", err)
	}

	return resp, nil
}

func (c *Client) do(ctx context.Context, client *http.Client, method string, camera Camera, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s%s", camera.Host, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/xml")
	}

	return client.Do(req)
}
