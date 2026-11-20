package isapi

import (
	"net/http"
	"sync"
	"time"

	"github.com/andriantp/camera"
)

type Camera struct {
	camera.Credential
	Channel int
}

type Client struct {
	http         *http.Client
	downloadHTTP *http.Client
	moveTimeout  time.Duration

	// auth caches digest clients per credential (credentialKey -> *authClients).
	auth sync.Map
}

func New(client *camera.Client) *Client {
	return &Client{
		http:         client.HTTP(),
		downloadHTTP: client.HTTPDownload(),
		moveTimeout:  client.MoveTimeout(),
	}
}
