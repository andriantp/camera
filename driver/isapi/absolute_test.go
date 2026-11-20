package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/andriantp/camera/model"
)

func TestPanAlmostEqual(t *testing.T) {
	tests := []struct {
		a, b float64
		want bool
	}{
		{0, 0, true},
		{0, 3600, true},
		{3600, 0, true},
		{3599.95, 0.02, true},
		{1800, 1800.05, true},
		{0, 0.2, false},
		{0, 1800, false},
		{10, 3590, false},
	}

	for _, tt := range tests {
		if got := panAlmostEqual(tt.a, tt.b); got != tt.want {
			t.Errorf("panAlmostEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestAlmostEqual(t *testing.T) {
	if !almostEqual(15, 1.5*10) {
		t.Error("almostEqual should tolerate float rounding")
	}
	if almostEqual(15, 15.2) {
		t.Error("almostEqual should reject differences above tolerance")
	}
}

// fakePTZ simulates a PTZ camera: PUT absoluteEx stores the requested
// position and GET absoluteEx reports it back.
type fakePTZ struct {
	mu                       sync.Mutex
	azimuth, elevation, zoom string
}

func (f *fakePTZ) handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/capabilities"):
		_, _ = w.Write([]byte(ptzCapabilitiesXML))

	case strings.HasSuffix(r.URL.Path, "/absoluteEx") && r.Method == http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.azimuth = between(string(body), "<azimuth>", "</azimuth>", f.azimuth)
		f.elevation = between(string(body), "<elevation>", "</elevation>", f.elevation)
		f.zoom = between(string(body), "<absoluteZoom>", "</absoluteZoom>", f.zoom)
		f.mu.Unlock()

	case strings.HasSuffix(r.URL.Path, "/absoluteEx"):
		f.mu.Lock()
		fmt.Fprintf(w, `<PTZAbsoluteEx><azimuth>%s</azimuth><elevation>%s</elevation><absoluteZoom>%s</absoluteZoom></PTZAbsoluteEx>`,
			f.azimuth, f.elevation, f.zoom)
		f.mu.Unlock()

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func between(s, start, end, fallback string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return fallback
	}
	s = s[i+len(start):]
	j := strings.Index(s, end)
	if j < 0 {
		return fallback
	}
	return s[:j]
}

func TestSetAbsolute(t *testing.T) {
	ptz := &fakePTZ{azimuth: "0", elevation: "0", zoom: "1"}
	client, cam := newTestCamera(t, ptz.handler)

	got, err := client.SetAbsolute(context.Background(), cam, model.PanTilt{X: 1800, Y: 450}, 15)
	if err != nil {
		t.Fatal(err)
	}

	if !almostEqual(got.Position.X, 1800) || !almostEqual(got.Position.Y, 450) || !almostEqual(got.Zoom, 15) {
		t.Errorf("got %+v, want x=1800 y=450 zoom=15", got)
	}
}

func TestSetAbsoluteOutOfRange(t *testing.T) {
	client, cam := newTestCamera(t, (&fakePTZ{}).handler)

	tests := []struct {
		name    string
		pos     model.PanTilt
		zoom    float64
		wantErr string
	}{
		{"pan", model.PanTilt{X: 3700, Y: 0}, 10, "pan out of range"},
		{"tilt", model.PanTilt{X: 0, Y: 2800}, 10, "tilt out of range"},
		{"zoom", model.PanTilt{X: 0, Y: 0}, 400, "zoom out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.SetAbsolute(context.Background(), cam, tt.pos, tt.zoom)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want %q", err, tt.wantErr)
			}
		})
	}
}
