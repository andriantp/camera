package isapi

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/andriantp/camera/model"
)

func TestRelativeMoveWaitsForDuration(t *testing.T) {
	var (
		mu              sync.Mutex
		movedAt, readAt time.Time
	)

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if r.Method == http.MethodPut {
			movedAt = time.Now()
			return
		}
		readAt = time.Now()
		_, _ = w.Write([]byte(statusXML))
	})

	const duration = 200 * time.Millisecond

	status, err := client.RelativeMove(context.Background(), cam, model.Speed{Pan: 40}, duration)
	if err != nil {
		t.Fatal(err)
	}
	if status.Position.X != 1800 {
		t.Errorf("pan = %v, want 1800", status.Position.X)
	}

	mu.Lock()
	defer mu.Unlock()
	if gap := readAt.Sub(movedAt); gap < duration {
		t.Errorf("status read %v after move, want >= %v", gap, duration)
	}
}

func TestRelativeMoveContextCanceledWhileWaiting(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.RelativeMove(ctx, cam, model.Speed{Pan: 40}, 5*time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("returned after %v, want prompt return", elapsed)
	}
}

func TestRelativeMoveRejected(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("no ptz"))
	})

	_, err := client.RelativeMove(context.Background(), cam, model.Speed{Pan: 40}, time.Millisecond)
	if err == nil {
		t.Fatal("expected error for rejected move")
	}
}
