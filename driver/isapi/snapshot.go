package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// maxErrorBody limits how much of an error response body is read into memory.
const maxErrorBody = 4 << 10

type SnapshotOptions struct {
	Width        int
	Height       int
	ImageQuality string
}

// GetSnapshot fetches a JPEG snapshot from the camera and returns the image
// as a stream.
//
// The returned size is the response Content-Length and is -1 when the camera
// does not send one (e.g. chunked encoding), so do not rely on it to size a
// buffer. The caller must read the body to EOF and close it; otherwise the
// underlying connection cannot be reused.
func (c *Client) GetSnapshot(ctx context.Context, camera Camera, opts SnapshotOptions) (int64, io.ReadCloser, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		snapshotPath(camera.Channel, opts),
		nil,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("get snapshot: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body := readErrorBody(resp)

		return 0, nil, fmt.Errorf(
			"get snapshot: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	// Some firmware answers 200 with an XML error instead of an image.
	if contentType := resp.Header.Get("Content-Type"); contentType != "" &&
		!strings.HasPrefix(strings.ToLower(contentType), "image/") {
		body := readErrorBody(resp)

		return 0, nil, fmt.Errorf(
			"get snapshot: unexpected content type %q: %s",
			contentType,
			body,
		)
	}

	return resp.ContentLength, resp.Body, nil
}

// readErrorBody reads at most maxErrorBody bytes of resp.Body for an error
// message, then closes the body.
func readErrorBody(resp *http.Response) []byte {
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	return body
}

func snapshotPath(channel int, opts SnapshotOptions) string {
	values := url.Values{}

	if opts.Width > 0 {
		values.Set("videoResolutionWidth", strconv.Itoa(opts.Width))
	}

	if opts.Height > 0 {
		values.Set("videoResolutionHeight", strconv.Itoa(opts.Height))
	}

	if opts.ImageQuality == "" {
		opts.ImageQuality = "better" //best, better, normal, general
	}

	values.Set("snapShotImageType", "JPEG")
	values.Set("imageQuality", opts.ImageQuality)

	return fmt.Sprintf(
		"/ISAPI/Streaming/channels/%d/picture?%s",
		channel,
		values.Encode(),
	)
}
