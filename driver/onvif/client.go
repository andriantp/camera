package onvif

import (
	"net/http"

	"github.com/andriantp/camera"
)

type Camera struct {
	camera.Credential
	Profile string
}

type Client struct {
	http *http.Client
}

func New(client *camera.Client) *Client {
	return &Client{
		http: client.HTTP(),
	}
}
