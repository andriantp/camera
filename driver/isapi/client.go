package isapi

import (
	"net/http"

	"github.com/andriantp/camera"
)

type Camera struct {
	camera.Credential
	Channel int
}

type Client struct {
	http         *http.Client
	downloadHTTP *http.Client
}

func New(client *camera.Client) *Client {
	return &Client{
		http:         client.HTTP(),
		downloadHTTP: client.HTTPDownload(),
	}
}
