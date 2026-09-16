package isapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetStreaming(ctx context.Context, camera Camera) ([]model.StreamingChannel, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/Streaming/channels",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseStreaming(resp)
}

func (c *Client) GetStreamingCapabilities(ctx context.Context, camera Camera) (*model.StreamingCapabilities, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf("/ISAPI/Streaming/channels/%d/capabilities", camera.Channel),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseStreamingCapabilities(resp)
}
