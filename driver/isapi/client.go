package isapi

import (
	"net/http"
	"sync"

	"github.com/andriantp/camera"
)

type Camera struct {
	camera.Credential
	Channel int
}

type Client struct {
	http         *http.Client
	downloadHTTP *http.Client

	// auth caches digest clients per credential (credentialKey -> *authClients).
	auth sync.Map
}

func New(client *camera.Client) *Client {
	return &Client{
		http:         client.HTTP(),
		downloadHTTP: client.HTTPDownload(),
	}
}
