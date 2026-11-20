package xml

import (
	"io"
	"net/http"
	"strings"
)

// response builds an *http.Response with the given status and body.
func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{},
	}
}

// trackingBody records whether Close was called.
type trackingBody struct {
	io.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

func trackedResponse(status int, body string) (*http.Response, *trackingBody) {
	tb := &trackingBody{Reader: strings.NewReader(body)}
	return &http.Response{StatusCode: status, Body: tb, Header: http.Header{}}, tb
}
