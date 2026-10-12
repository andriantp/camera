package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetMotionDetection(ctx context.Context, camera Camera) (*model.MotionDetection, error) {
	path := fmt.Sprintf(
		"/ISAPI/System/Video/inputs/channels/%d/motionDetection",
		camera.Channel,
	)

	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"get motion detection: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return xml.ParseMotionDetection(resp)
}
