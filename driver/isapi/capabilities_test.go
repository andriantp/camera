package isapi

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGetCapabilitiesPartial(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		if !requireDigest(w, r) {
			return
		}
		if strings.Contains(r.URL.Path, "/PTZCtrl/") {
			_, _ = w.Write([]byte(ptzCapabilitiesXML))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	caps, err := client.GetCapabilities(context.Background(), cam)
	if err != nil {
		t.Fatal(err)
	}

	if caps.Isapi.PTZ == nil || caps.Isapi.PTZ.Preset.MaxPreset != 256 {
		t.Errorf("PTZ = %+v, want parsed capabilities", caps.Isapi.PTZ)
	}
	if caps.Isapi.System != nil || caps.Isapi.Streaming != nil || caps.Isapi.Analytic != nil {
		t.Error("unsupported groups should be nil")
	}
}

func TestGetCapabilitiesAllFail(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	caps, err := client.GetCapabilities(context.Background(), cam)
	if err == nil || caps != nil {
		t.Fatalf("caps = %v, err = %v; want nil caps and an error", caps, err)
	}

	for _, group := range []string{"system", "streaming", "ptz", "event", "image", "analytic"} {
		if !strings.Contains(err.Error(), group+":") {
			t.Errorf("error does not mention %q: %v", group, err)
		}
	}
}

func TestGetCapabilitiesContextCanceled(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := client.GetCapabilities(ctx, cam); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context canceled", err)
	}
}

func TestGetCapabilitiesConcurrent(t *testing.T) {
	const delay = 100 * time.Millisecond

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusNotFound)
	})

	start := time.Now()
	_, _ = client.GetCapabilities(context.Background(), cam)

	// Six groups fetched sequentially would take ~6x delay.
	if elapsed := time.Since(start); elapsed > 4*delay {
		t.Errorf("took %v, want groups fetched concurrently (< %v)", elapsed, 4*delay)
	}
}
