package isapi

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

const statusXML = `<PTZStatus><AbsoluteHigh><azimuth>1800</azimuth><elevation>100</elevation><absoluteZoom>10</absoluteZoom></AbsoluteHigh></PTZStatus>`

func TestRequestReusesDigestChallenge(t *testing.T) {
	var total, challenged atomic.Int32

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		total.Add(1)
		if !requireDigest(w, r) {
			challenged.Add(1)
			return
		}
		_, _ = w.Write([]byte(statusXML))
	})

	const calls = 5
	for i := 0; i < calls; i++ {
		if _, err := client.GetStatus(context.Background(), cam); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}

	// One 401 challenge for the first call, then one request per call.
	if got := challenged.Load(); got != 1 {
		t.Errorf("challenges = %d, want 1", got)
	}
	if got := total.Load(); got != calls+1 {
		t.Errorf("requests = %d, want %d", got, calls+1)
	}
}

func TestRequestDigestConcurrent(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		if requireDigest(w, r) {
			_, _ = w.Write([]byte(statusXML))
		}
	})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetStatus(context.Background(), cam); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestRequestSeparateCredentialsGetSeparateClients(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {})

	other := cam
	other.Password = "different"

	first := client.authClientsFor(cam)
	again := client.authClientsFor(cam)
	different := client.authClientsFor(other)

	if first != again {
		t.Fatal("same credentials do not reuse the digest client")
	}
	if first == different {
		t.Fatal("different credentials share the same digest client")
	}
}

func TestRequestContentTypeOnlyWithBody(t *testing.T) {
	contentTypes := map[string]string{}
	var mu sync.Mutex

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		if !requireDigest(w, r) {
			return
		}
		mu.Lock()
		contentTypes[r.Method] = r.Header.Get("Content-Type")
		mu.Unlock()
	})

	ctx := context.Background()

	resp, err := client.request(ctx, http.MethodGet, cam, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	resp, err = client.request(ctx, http.MethodPut, cam, "/x", strings.NewReader("<a/>"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got := contentTypes[http.MethodGet]; got != "" {
		t.Errorf("GET Content-Type = %q, want empty", got)
	}
	if got := contentTypes[http.MethodPut]; got != "application/xml" {
		t.Errorf("PUT Content-Type = %q, want application/xml", got)
	}
}
