package isapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andriantp/camera"
	"github.com/andriantp/camera/model"
)

// stuckPTZ accepts move commands but never changes its reported position.
func stuckPTZ(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		_, _ = io.Copy(io.Discard, r.Body)
		return
	}
	(&fakePTZ{azimuth: "0", elevation: "0", zoom: "1"}).handler(w, r)
}

func newStuckCamera(t *testing.T, moveTimeout time.Duration) (*Client, Camera) {
	t.Helper()

	_, cam := newTestCamera(t, stuckPTZ)
	client := New(camera.New(camera.Config{MoveTimeout: moveTimeout}))

	return client, cam
}

func TestMoveTimeoutIsConfigurable(t *testing.T) {
	const moveTimeout = 300 * time.Millisecond

	client, cam := newStuckCamera(t, moveTimeout)

	tests := map[string]func() error{
		"SetZoom": func() error {
			_, err := client.SetZoom(context.Background(), cam, 20)
			return err
		},
		"SetAbsolute": func() error {
			_, err := client.SetAbsolute(context.Background(), cam, model.PanTilt{X: 900, Y: 0}, 10)
			return err
		},
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			start := time.Now()
			err := call()
			elapsed := time.Since(start)

			if !errors.Is(err, ErrMoveTimeout) {
				t.Fatalf("err = %v, want ErrMoveTimeout", err)
			}
			if !strings.Contains(err.Error(), "300ms") {
				t.Errorf("error should mention the timeout: %v", err)
			}
			if elapsed < moveTimeout || elapsed > moveTimeout+time.Second {
				t.Errorf("gave up after %v, want about %v", elapsed, moveTimeout)
			}
		})
	}
}

func TestMoveTimeoutContextDeadlineWins(t *testing.T) {
	client, cam := newStuckCamera(t, time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.SetZoom(ctx, cam, 20)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context deadline exceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("returned after %v, want the shorter context deadline to apply", elapsed)
	}
}

func TestMoveTimeoutDefault(t *testing.T) {
	client := New(camera.New(camera.Config{}))

	if client.moveTimeout != defaultMoveTimeout {
		t.Errorf("moveTimeout = %v, want %v", client.moveTimeout, defaultMoveTimeout)
	}
}

func TestWaitPositionGuardsZeroTimeout(t *testing.T) {
	// A Client built without camera.New (zero moveTimeout) must not time out
	// immediately; it falls back to the default.
	_, cam := newTestCamera(t, (&fakePTZ{azimuth: "0", elevation: "0", zoom: "1"}).handler)
	client := New(camera.New(camera.Config{}))
	client.moveTimeout = 0

	got, err := client.waitPosition(context.Background(), cam, func(*model.Absolute) bool { return true }, "any")
	if err != nil || got == nil {
		t.Fatalf("got %v, %v", got, err)
	}
}
