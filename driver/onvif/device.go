package onvif

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"

	usego "github.com/use-go/onvif"
)

type deviceKey struct {
	host     string
	username string
	password string
}

// cachedDevice is a use-go/onvif Device reused across calls to one camera.
//
// usego.NewDevice performs a GetCapabilities round trip to discover service
// endpoints, so creating a Device per call doubles the SOAP traffic. The
// Device's HTTP client cannot be replaced after creation, so to keep
// per-call context support its transport reads the context of the call that
// currently holds mu. Calls to the same camera are therefore serialized.
type cachedDevice struct {
	mu        sync.Mutex
	dev       *usego.Device
	transport *contextTransport
}

// withDevice runs fn with the camera's cached Device, binding all HTTP
// requests made inside fn to ctx.
func (c *Client) withDevice(ctx context.Context, camera Camera, fn func(dev *usego.Device) error) error {
	cd := c.cachedDevice(camera)

	cd.mu.Lock()
	defer cd.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	cd.transport.ctx = ctx
	defer func() { cd.transport.ctx = context.Background() }()

	if cd.dev == nil {
		dev, err := usego.NewDevice(usego.DeviceParams{
			Xaddr:      onvifAddr(camera.Host),
			Username:   camera.Username,
			Password:   camera.Password,
			HttpClient: cd.httpClient(c.http),
		})
		if err != nil {
			// NewDevice replaces the underlying error with a generic message;
			// surface cancellation/deadline errors so callers can detect them.
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			return err
		}
		cd.dev = dev
	}

	return fn(cd.dev)
}

func (c *Client) cachedDevice(camera Camera) *cachedDevice {
	key := deviceKey{host: camera.Host, username: camera.Username, password: camera.Password}

	if v, ok := c.devices.Load(key); ok {
		return v.(*cachedDevice)
	}

	base := c.http.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	v, _ := c.devices.LoadOrStore(key, &cachedDevice{
		transport: &contextTransport{ctx: context.Background(), base: base},
	})

	return v.(*cachedDevice)
}

func (cd *cachedDevice) httpClient(base *http.Client) *http.Client {
	client := *base
	client.Transport = cd.transport
	return &client
}

// contextTransport attaches ctx to every outgoing request, because
// use-go/onvif builds its requests without a context. ctx is only written
// while the owning cachedDevice's mutex is held.
type contextTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.base.RoundTrip(req.WithContext(t.ctx))
}

// onvifAddr converts a camera Host into the "host[:port]" that use-go/onvif
// expects; the library always prepends "http://" itself, so ONVIF in this
// package runs over HTTP only.
//
//	"192.168.1.64:8080"         -> "192.168.1.64:8080"
//	"http://192.168.1.64:8080"  -> "192.168.1.64:8080"
//	"https://192.168.1.64:8443" -> "192.168.1.64"
//
// For https:// the port belongs to HTTPS, so it is dropped and ONVIF uses the
// default HTTP port.
//
// A Credential shared with the ISAPI driver can therefore use an https://
// host. If ONVIF listens on a non-default HTTP port, use a separate
// Credential with "host:port" for the ONVIF client.
func onvifAddr(host string) string {
	scheme, rest, found := strings.Cut(host, "://")
	if !found {
		return strings.TrimRight(host, "/")
	}

	rest, _, _ = strings.Cut(rest, "/")

	if strings.EqualFold(scheme, "https") {
		if h, _, err := net.SplitHostPort(rest); err == nil {
			if strings.Contains(h, ":") { // IPv6 literal needs brackets in a URL
				return "[" + h + "]"
			}
			return h
		}
	}

	return rest
}
