package isapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetSnapshot(t *testing.T) {
	const jpeg = "\xff\xd8\xff\xe0fake-jpeg"

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte(jpeg))
	})

	size, body, err := client.GetSnapshot(context.Background(), cam, SnapshotOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if string(data) != jpeg {
		t.Errorf("body = %q, want %q", data, jpeg)
	}
	if size != int64(len(jpeg)) {
		t.Errorf("size = %d, want %d", size, len(jpeg))
	}
}

func TestGetSnapshotErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr string
	}{
		{
			name: "xml with 200",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				_, _ = w.Write([]byte(`<ResponseStatus><statusString>Device Busy</statusString></ResponseStatus>`))
			},
			wantErr: "Device Busy",
		},
		{
			name: "non-200",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte("forbidden"))
			},
			wantErr: "status 403: forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cam := newTestCamera(t, tt.handler)

			_, _, err := client.GetSnapshot(context.Background(), cam, SnapshotOptions{})
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestGetSnapshotLimitsErrorBody(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(strings.Repeat("x", 1<<20)))
	})

	_, _, err := client.GetSnapshot(context.Background(), cam, SnapshotOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(err.Error()) > maxErrorBody+100 {
		t.Errorf("error message is %d bytes, want it capped near %d", len(err.Error()), maxErrorBody)
	}
}

func TestSnapshotPath(t *testing.T) {
	path := snapshotPath(2, SnapshotOptions{Width: 1920, Height: 1080})

	u, err := url.Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	if u.Path != "/ISAPI/Streaming/channels/2/picture" {
		t.Errorf("path = %q", u.Path)
	}

	q := u.Query()
	want := map[string]string{
		"videoResolutionWidth":  "1920",
		"videoResolutionHeight": "1080",
		"snapShotImageType":     "JPEG",
		"imageQuality":          "better",
	}
	for k, v := range want {
		if got := q.Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}

	if q := mustQuery(t, snapshotPath(1, SnapshotOptions{})); q.Has("videoResolutionWidth") {
		t.Error("width should be omitted when zero")
	}
}

func mustQuery(t *testing.T, path string) url.Values {
	t.Helper()
	u, err := url.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	return u.Query()
}
