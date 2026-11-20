package isapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/andriantp/camera"
)

func TestBaseURL(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{"192.168.1.64", "http://192.168.1.64"},
		{"192.168.1.64:8080", "http://192.168.1.64:8080"},
		{"camera.local", "http://camera.local"},
		{"http://192.168.1.64", "http://192.168.1.64"},
		{"https://192.168.1.64", "https://192.168.1.64"},
		{"https://192.168.1.64:8443", "https://192.168.1.64:8443"},
		{"https://192.168.1.64/", "https://192.168.1.64"},
		{"https://192.168.1.64:8443/ignored/path", "https://192.168.1.64:8443"},
		{"[fe80::1]:8080", "http://[fe80::1]:8080"},
		{"https://[fe80::1]:8443", "https://[fe80::1]:8443"},
	}

	for _, tt := range tests {
		got, err := baseURL(tt.host)
		if err != nil {
			t.Errorf("baseURL(%q) error: %v", tt.host, err)
			continue
		}
		if got != tt.want {
			t.Errorf("baseURL(%q) = %q, want %q", tt.host, got, tt.want)
		}
	}
}

func TestBaseURLInvalid(t *testing.T) {
	for _, host := range []string{"", "htp://192.168.1.64", "ftp://192.168.1.64", "https://", "http://bad host"} {
		if got, err := baseURL(host); err == nil {
			t.Errorf("baseURL(%q) = %q, want error", host, got)
		}
	}
}

func TestRequestInvalidHost(t *testing.T) {
	client := New(camera.New(camera.Config{}))
	cam := Camera{Credential: camera.Credential{Host: "ftp://192.168.1.64"}, Channel: 1}

	_, err := client.GetStatus(context.Background(), cam)
	if err == nil || !strings.Contains(err.Error(), "scheme must be http or https") {
		t.Fatalf("err = %v", err)
	}
}

// newTLSCamera starts an HTTPS fake camera with a self-signed certificate
// that serves PTZ status behind digest auth.
func newTLSCamera(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var hits atomic.Int32
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			t.Error("request did not arrive over TLS")
		}
		if !requireDigest(w, r) {
			return
		}
		hits.Add(1)
		_, _ = w.Write([]byte(statusXML))
	}))
	t.Cleanup(srv.Close)

	return srv, &hits
}

func tlsCamera(srv *httptest.Server) Camera {
	return Camera{
		Credential: camera.Credential{Host: srv.URL, Username: "admin", Password: "secret"},
		Channel:    1,
	}
}

func TestHTTPSRejectsUntrustedCertificateByDefault(t *testing.T) {
	srv, hits := newTLSCamera(t)
	client := New(camera.New(camera.Config{}))

	_, err := client.GetStatus(context.Background(), tlsCamera(srv))

	var certErr *tls.CertificateVerificationError
	if !errors.As(err, &certErr) {
		t.Fatalf("err = %v, want certificate verification error", err)
	}
	if hits.Load() != 0 {
		t.Error("request reached the camera despite an untrusted certificate")
	}
}

func TestHTTPSWithTrustedRootCA(t *testing.T) {
	srv, hits := newTLSCamera(t)

	roots := x509.NewCertPool()
	roots.AddCert(srv.Certificate())

	client := New(camera.New(camera.Config{
		TLSConfig: &tls.Config{RootCAs: roots},
	}))

	status, err := client.GetStatus(context.Background(), tlsCamera(srv))
	if err != nil {
		t.Fatal(err)
	}
	if status.Position.X != 1800 || hits.Load() != 1 {
		t.Errorf("status = %+v, hits = %d", status, hits.Load())
	}
}

func TestHTTPSWithInsecureSkipVerify(t *testing.T) {
	srv, _ := newTLSCamera(t)

	client := New(camera.New(camera.Config{
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // test against self-signed server
	}))

	if _, err := client.GetStatus(context.Background(), tlsCamera(srv)); err != nil {
		t.Fatal(err)
	}
}

func TestHTTPSTLSConfigIsCloned(t *testing.T) {
	srv, _ := newTLSCamera(t)

	cfg := &tls.Config{InsecureSkipVerify: true} //nolint:gosec // test against self-signed server
	client := New(camera.New(camera.Config{TLSConfig: cfg}))

	// Mutating the caller's config after New must not affect the client.
	cfg.InsecureSkipVerify = false

	if _, err := client.GetStatus(context.Background(), tlsCamera(srv)); err != nil {
		t.Fatalf("client was affected by later TLSConfig mutation: %v", err)
	}
}

func TestHTTPSContentDownloadUsesTLS(t *testing.T) {
	var overTLS atomic.Bool

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireDigest(w, r) {
			return
		}
		overTLS.Store(r.TLS != nil)
		_, _ = w.Write([]byte("video"))
	}))
	t.Cleanup(srv.Close)

	client := New(camera.New(camera.Config{
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // test against self-signed server
	}))

	_, body, err := client.ContentDownload(context.Background(), tlsCamera(srv), "rtsp://cam/track?size=5")
	if err != nil {
		t.Fatal(err)
	}
	body.Close()

	if !overTLS.Load() {
		t.Error("download (which carries the password in its body) was not sent over TLS")
	}
}
