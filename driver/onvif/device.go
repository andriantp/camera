package onvif

import (
	"context"
	"net/http"

	usego "github.com/use-go/onvif"
)

// device creates a use-go/onvif Device whose HTTP requests are bound to ctx.
// use-go/onvif does not accept a context, so ctx is injected at the transport
// level; cancelling ctx aborts any in-flight SOAP call.
func (c *Client) device(ctx context.Context, camera Camera) (*usego.Device, error) {
	dev, err := usego.NewDevice(usego.DeviceParams{
		Xaddr:      camera.Host,
		Username:   camera.Username,
		Password:   camera.Password,
		HttpClient: c.httpClientWithContext(ctx),
	})
	if err != nil {
		// NewDevice replaces the underlying error with a generic message;
		// surface cancellation/deadline errors so callers can detect them.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, err
	}

	return dev, nil
}

func (c *Client) httpClientWithContext(ctx context.Context) *http.Client {
	base := c.http.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	client := *c.http
	client.Transport = &contextTransport{ctx: ctx, base: base}

	return &client
}

type contextTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.base.RoundTrip(req.WithContext(t.ctx))
}
