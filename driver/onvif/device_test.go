package onvif

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andriantp/camera"
	"github.com/andriantp/camera/model"
)

// fakeCamera is a minimal ONVIF SOAP endpoint.
type fakeCamera struct {
	srv *httptest.Server

	capabilitiesCalls atomic.Int32
	statusCalls       atomic.Int32
	moveCalls         atomic.Int32

	unavailable atomic.Bool  // fail GetCapabilities (device discovery)
	rejectMove  atomic.Bool  // answer RelativeMove with a SOAP fault
	statusDelay atomic.Int64 // delay (ns) before answering GetStatus
}

func newFakeCamera(t *testing.T) *fakeCamera {
	t.Helper()

	f := &fakeCamera{}
	f.srv = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.srv.Close)

	return f
}

func (f *fakeCamera) handle(w http.ResponseWriter, r *http.Request) {
	// Read the body first so the server notices client disconnects.
	body, _ := io.ReadAll(r.Body)
	soap := string(body)

	switch {
	case strings.Contains(soap, "GetCapabilities"):
		f.capabilitiesCalls.Add(1)
		if f.unavailable.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `<Envelope><Body><GetCapabilitiesResponse><Capabilities>
			<Media><XAddr>%[1]s/onvif/media</XAddr></Media>
			<PTZ><XAddr>%[1]s/onvif/ptz</XAddr></PTZ>
		</Capabilities></GetCapabilitiesResponse></Body></Envelope>`, f.srv.URL)

	case strings.Contains(soap, "GetStatus"):
		f.statusCalls.Add(1)
		select {
		case <-time.After(time.Duration(f.statusDelay.Load())):
		case <-r.Context().Done():
			return
		}
		_, _ = w.Write([]byte(`<Envelope><Body><GetStatusResponse><PTZStatus>
			<Position><PanTilt x="0.5" y="-0.25"/></Position>
			<MoveStatus><PanTilt>IDLE</PanTilt></MoveStatus>
		</PTZStatus></GetStatusResponse></Body></Envelope>`))

	case strings.Contains(soap, "RelativeMove"):
		f.moveCalls.Add(1)
		if f.rejectMove.Load() {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("SOAP fault"))
			return
		}
		_, _ = w.Write([]byte(`<Envelope><Body><RelativeMoveResponse/></Body></Envelope>`))

	case strings.Contains(soap, "GetConfigurations"):
		_, _ = w.Write([]byte(`<Envelope><Body><GetConfigurationsResponse><PTZConfiguration>
			<DefaultPTZSpeed><PanTilt x="0.5" y="0.5"/></DefaultPTZSpeed>
			<PanTiltLimits><Range><XRange><Min>-1</Min><Max>1</Max></XRange><YRange><Min>-1</Min><Max>1</Max></YRange></Range></PanTiltLimits>
		</PTZConfiguration></GetConfigurationsResponse></Body></Envelope>`))

	case strings.Contains(soap, "GetProfiles"):
		_, _ = w.Write([]byte(`<Envelope><Body><GetProfilesResponse>
			<Profiles token="profile_1"><Name>main</Name></Profiles>
		</GetProfilesResponse></Body></Envelope>`))

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (f *fakeCamera) camera() Camera {
	return Camera{
		Credential: camera.Credential{
			Host:     strings.TrimPrefix(f.srv.URL, "http://"),
			Username: "admin",
			Password: "secret",
		},
		Profile: "profile_1",
	}
}

func newClient() *Client {
	return New(camera.New(camera.Config{}))
}

func TestDeviceIsCachedAcrossCalls(t *testing.T) {
	cam := newFakeCamera(t)
	client := newClient()

	for i := 0; i < 5; i++ {
		status, err := client.GetStatus(context.Background(), cam.camera())
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if status.Position.X != 0.5 || status.Position.Y != -0.25 || status.MoveStatus != "IDLE" {
			t.Fatalf("status = %+v", status)
		}
	}

	if got := cam.capabilitiesCalls.Load(); got != 1 {
		t.Errorf("GetCapabilities calls = %d, want 1 (device should be cached)", got)
	}
}

func TestDeviceConcurrentCalls(t *testing.T) {
	cam := newFakeCamera(t)
	cam.statusDelay.Store(int64(5 * time.Millisecond))
	client := newClient()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if got := cam.capabilitiesCalls.Load(); got != 1 {
		t.Errorf("GetCapabilities calls = %d, want 1", got)
	}
	if got := cam.statusCalls.Load(); got != 20 {
		t.Errorf("GetStatus calls = %d, want 20", got)
	}
}

func TestDeviceSeparateCamerasNotShared(t *testing.T) {
	camA, camB := newFakeCamera(t), newFakeCamera(t)
	client := newClient()

	for _, cam := range []*fakeCamera{camA, camB} {
		if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
			t.Fatal(err)
		}
	}

	if camA.capabilitiesCalls.Load() != 1 || camB.capabilitiesCalls.Load() != 1 {
		t.Error("each camera should get its own device discovery")
	}
}

func TestDeviceContextDeadline(t *testing.T) {
	cam := newFakeCamera(t)
	cam.statusDelay.Store(int64(5 * time.Second))
	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.GetStatus(ctx, cam.camera())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("returned after %v, want prompt return", elapsed)
	}

	// The cached device must still work with a fresh context.
	cam.statusDelay.Store(0)
	if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
		t.Fatalf("after deadline: %v", err)
	}
}

func TestDeviceContextDeadlineDuringDiscovery(t *testing.T) {
	// A server that never answers device discovery.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		<-r.Context().Done()
	}))
	t.Cleanup(srv.Close)

	cam := Camera{Credential: camera.Credential{Host: strings.TrimPrefix(srv.URL, "http://")}}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if _, err := newClient().GetStatus(ctx, cam); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
}

func TestDeviceFailedDiscoveryIsRetried(t *testing.T) {
	cam := newFakeCamera(t)
	cam.unavailable.Store(true)
	client := newClient()

	if _, err := client.GetStatus(context.Background(), cam.camera()); err == nil {
		t.Fatal("expected error while camera is unavailable")
	}

	cam.unavailable.Store(false)
	if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
		t.Fatalf("expected recovery once camera is available: %v", err)
	}

	if got := cam.capabilitiesCalls.Load(); got != 2 {
		t.Errorf("GetCapabilities calls = %d, want 2", got)
	}
}

func TestRelativeMove(t *testing.T) {
	cam := newFakeCamera(t)
	client := newClient()

	start := time.Now()
	status, err := client.RelativeMove(context.Background(), cam.camera(), model.PanTilt{X: 0.1}, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	if time.Since(start) < 50*time.Millisecond {
		t.Error("RelativeMove returned before the wait elapsed")
	}
	if status.Position.X != 0.5 {
		t.Errorf("status = %+v", status)
	}
	if cam.moveCalls.Load() != 1 || cam.statusCalls.Load() != 1 {
		t.Errorf("move=%d status=%d, want 1 each", cam.moveCalls.Load(), cam.statusCalls.Load())
	}
}

func TestRelativeMoveDoesNotHoldLockWhileWaiting(t *testing.T) {
	cam := newFakeCamera(t)
	client := newClient()

	// Warm the device cache.
	if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = client.RelativeMove(context.Background(), cam.camera(), model.PanTilt{X: 0.1}, time.Second)
	}()

	// Give RelativeMove time to send the move and start waiting.
	time.Sleep(100 * time.Millisecond)

	start := time.Now()
	if _, err := client.GetStatus(context.Background(), cam.camera()); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("GetStatus blocked for %v while RelativeMove was waiting", elapsed)
	}

	<-done
}
