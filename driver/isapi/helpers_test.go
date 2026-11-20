package isapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andriantp/camera"
)

// newTestCamera starts a fake camera serving h and returns a client and a
// Camera pointing at it. The server is closed when the test ends.
func newTestCamera(t *testing.T, h http.HandlerFunc) (*Client, Camera) {
	t.Helper()

	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	client := New(camera.New(camera.Config{}))
	cam := Camera{
		Credential: camera.Credential{
			Host:     strings.TrimPrefix(srv.URL, "http://"),
			Username: "admin",
			Password: "secret",
		},
		Channel: 1,
	}

	return client, cam
}

// requireDigest drains the request body and answers 401 with a digest
// challenge when the request carries no digest Authorization header.
// It reports whether the request was authenticated.
func requireDigest(w http.ResponseWriter, r *http.Request) bool {
	_, _ = io.Copy(io.Discard, r.Body)

	if strings.HasPrefix(r.Header.Get("Authorization"), "Digest ") {
		return true
	}

	w.Header().Set("WWW-Authenticate", `Digest realm="camera", nonce="abc123", qop="auth", algorithm=MD5`)
	w.WriteHeader(http.StatusUnauthorized)

	return false
}

const ptzCapabilitiesXML = `<PTZChanelCap>
	<AbsolutePanTiltPositionSpace>
		<XRange><Min>0</Min><Max>3600</Max></XRange>
		<YRange><Min>-900</Min><Max>2700</Max></YRange>
	</AbsolutePanTiltPositionSpace>
	<AbsoluteZoomPositionSpace>
		<ZRange><Min>10</Min><Max>320</Max></ZRange>
	</AbsoluteZoomPositionSpace>
	<maxPresetNum>256</maxPresetNum>
</PTZChanelCap>`
