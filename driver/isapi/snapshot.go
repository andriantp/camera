package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type SnapshotOptions struct {
	Width        int
	Height       int
	ImageQuality string
}

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
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		return 0, nil, fmt.Errorf(
			"get snapshot: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return resp.ContentLength, resp.Body, nil
}

func snapshotPath(channel int, opts SnapshotOptions) string {
	if opts.Width == 0 {
		opts.Width = 1280
	}

	if opts.Height == 0 {
		opts.Height = 720
	}

	if opts.ImageQuality == "" {
		opts.ImageQuality = "better" //best, better, normal, general
	}

	values := url.Values{}

	values.Set("snapShotImageType", "JPEG")
	values.Set("videoResolutionWidth", strconv.Itoa(opts.Width))
	values.Set("videoResolutionHeight", strconv.Itoa(opts.Height))
	values.Set("imageQuality", opts.ImageQuality)

	return fmt.Sprintf(
		"/ISAPI/Streaming/channels/%d/picture?%s",
		channel,
		values.Encode(),
	)
}
