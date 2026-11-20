package isapi

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestSetZoom(t *testing.T) {
	ptz := &fakePTZ{azimuth: "0", elevation: "0", zoom: "1"}
	client, cam := newTestCamera(t, ptz.handler)

	// 1.5 * 10 is not exactly 15 in floating point; SetZoom must still
	// recognize that the target was reached instead of timing out.
	got, err := client.SetZoom(context.Background(), cam, 15)
	if err != nil {
		t.Fatal(err)
	}

	if !almostEqual(got.Zoom, 15) {
		t.Errorf("zoom = %v, want 15", got.Zoom)
	}
}

func TestSetZoomUnsupported(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<PTZChanelCap></PTZChanelCap>`))
	})

	_, err := client.SetZoom(context.Background(), cam, 10)
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("err = %v, want not supported", err)
	}
}

func TestSetZoomOutOfRange(t *testing.T) {
	client, cam := newTestCamera(t, (&fakePTZ{}).handler)

	_, err := client.SetZoom(context.Background(), cam, 500)
	if err == nil || !strings.Contains(err.Error(), "zoom out of range") {
		t.Fatalf("err = %v, want out of range", err)
	}
}
