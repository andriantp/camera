package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/icholy/digest"
)

func (c *Client) request(ctx context.Context, method string, camera Camera, path string, body io.Reader) (*http.Response, error) {

	url := fmt.Sprintf(
		"http://%s%s",
		camera.Host,
		path,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{
		Timeout: c.http.Timeout,
		Transport: &digest.Transport{
			Username:  camera.Username,
			Password:  camera.Password,
			Transport: c.http.Transport,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) requestDownload(ctx context.Context,method string,camera Camera,path string,body io.Reader) (*http.Response, error) {

	url := fmt.Sprintf("http://%s%s", camera.Host, path)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{
		Timeout: c.downloadHTTP.Timeout,
		Transport: &digest.Transport{
			Username:  camera.Username,
			Password:  camera.Password,
			Transport: c.downloadHTTP.Transport,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do download request: %w", err)
	}

	return resp, nil
}