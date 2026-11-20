package onvif

import (
	"net/http"
	"sync"

	"github.com/andriantp/camera"
)

type Camera struct {
	camera.Credential
	Profile string
}

type Client struct {
	http *http.Client

	// devices caches ONVIF devices per camera (deviceKey -> *cachedDevice).
	devices sync.Map
}

func New(client *camera.Client) *Client {
	return &Client{
		http: client.HTTP(),
	}
}
