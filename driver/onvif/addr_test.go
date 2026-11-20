package onvif

import (
	"context"
	"testing"
)

func TestOnvifAddr(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{"192.168.1.64", "192.168.1.64"},
		{"192.168.1.64:8080", "192.168.1.64:8080"},
		{"192.168.1.64/", "192.168.1.64"},
		{"http://192.168.1.64", "192.168.1.64"},
		{"http://192.168.1.64:8080", "192.168.1.64:8080"},
		{"http://192.168.1.64:8080/", "192.168.1.64:8080"},
		{"https://192.168.1.64", "192.168.1.64"},
		{"https://192.168.1.64:8443", "192.168.1.64"},
		{"HTTPS://192.168.1.64:8443", "192.168.1.64"},
		{"https://192.168.1.64:8443/path", "192.168.1.64"},
		{"https://[fe80::1]:8443", "[fe80::1]"},
		{"http://[fe80::1]:8080", "[fe80::1]:8080"},
	}

	for _, tt := range tests {
		if got := onvifAddr(tt.host); got != tt.want {
			t.Errorf("onvifAddr(%q) = %q, want %q", tt.host, got, tt.want)
		}
	}
}

func TestHostWithHTTPScheme(t *testing.T) {
	cam := newFakeCamera(t)

	c := cam.camera()
	c.Host = cam.srv.URL // "http://127.0.0.1:port"

	status, err := newClient().GetStatus(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if status.Position.X != 0.5 {
		t.Errorf("status = %+v", status)
	}
}
