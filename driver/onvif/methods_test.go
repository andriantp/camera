package onvif

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/andriantp/camera/model"
)

func TestGetCapabilities(t *testing.T) {
	cam := newFakeCamera(t)

	got, err := newClient().GetCapabilities(context.Background(), cam.camera())
	if err != nil {
		t.Fatal(err)
	}

	if got.Onvif == nil || !strings.HasSuffix(got.Onvif.PTZ.XAddr, "/onvif/ptz") || !strings.HasSuffix(got.Onvif.Media.XAddr, "/onvif/media") {
		t.Errorf("got %+v", got.Onvif)
	}
}

func TestGetConfigurations(t *testing.T) {
	cam := newFakeCamera(t)

	got, err := newClient().GetConfigurations(context.Background(), cam.camera())
	if err != nil {
		t.Fatal(err)
	}

	if got.DefaultPTZSpeed.PanTilt.X != 0.5 || got.PanTiltLimits.Range.XRange.Max != 1 {
		t.Errorf("got %+v", got)
	}
}

func TestGetProfiles(t *testing.T) {
	cam := newFakeCamera(t)

	got, err := newClient().GetProfiles(context.Background(), cam.camera())
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].Token != "profile_1" || got[0].Name != "main" {
		t.Errorf("got %+v", got)
	}
}

func TestMethodsFailWhenCameraUnavailable(t *testing.T) {
	cam := newFakeCamera(t)
	cam.unavailable.Store(true)
	client := newClient()
	ctx := context.Background()

	calls := map[string]func() error{
		"GetCapabilities":   func() error { _, err := client.GetCapabilities(ctx, cam.camera()); return err },
		"GetConfigurations": func() error { _, err := client.GetConfigurations(ctx, cam.camera()); return err },
		"GetProfiles":       func() error { _, err := client.GetProfiles(ctx, cam.camera()); return err },
		"GetStatus":         func() error { _, err := client.GetStatus(ctx, cam.camera()); return err },
		"RelativeMove": func() error {
			_, err := client.RelativeMove(ctx, cam.camera(), model.PanTilt{X: 0.1}, 0)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestMethodsFailWhenEndpointUnreachable(t *testing.T) {
	// Discovery succeeds, but every later SOAP call hits a closed port.
	cam := newFakeCamera(t)
	client := newClient()
	ctx := context.Background()

	if _, err := client.GetStatus(ctx, cam.camera()); err != nil {
		t.Fatal(err)
	}
	cam.srv.Close()

	calls := map[string]func() error{
		"GetCapabilities":   func() error { _, err := client.GetCapabilities(ctx, cam.camera()); return err },
		"GetConfigurations": func() error { _, err := client.GetConfigurations(ctx, cam.camera()); return err },
		"GetProfiles":       func() error { _, err := client.GetProfiles(ctx, cam.camera()); return err },
		"GetStatus":         func() error { _, err := client.GetStatus(ctx, cam.camera()); return err },
		"RelativeMove": func() error {
			_, err := client.RelativeMove(ctx, cam.camera(), model.PanTilt{X: 0.1}, 0)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestRelativeMoveRejected(t *testing.T) {
	cam := newFakeCamera(t)
	cam.rejectMove.Store(true)

	_, err := newClient().RelativeMove(context.Background(), cam.camera(), model.PanTilt{X: 0.1}, 0)
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Fatalf("err = %v, want status 400", err)
	}
	if cam.statusCalls.Load() != 0 {
		t.Error("status should not be read after a rejected move")
	}
}

func TestRelativeMoveContextCanceledWhileWaiting(t *testing.T) {
	cam := newFakeCamera(t)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := newClient().RelativeMove(ctx, cam.camera(), model.PanTilt{X: 0.1}, 5*time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	if cam.statusCalls.Load() != 0 {
		t.Error("status should not be read after cancellation")
	}
}

func TestWithDeviceRejectsCanceledContext(t *testing.T) {
	cam := newFakeCamera(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := newClient().GetStatus(ctx, cam.camera()); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context canceled", err)
	}
	if cam.capabilitiesCalls.Load() != 0 {
		t.Error("canceled context should not reach the camera")
	}
}
