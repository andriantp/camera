package isapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetPTZCapabilities(ctx context.Context, camera Camera) (*model.PTZCapabilities, error) {
	return c.getPTZCapabilities(ctx, camera)
}

func (c *Client) getPTZCapabilities(ctx context.Context, camera Camera) (*model.PTZCapabilities, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/capabilities",
			camera.Channel,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParsePTZCapabilities(resp)
}
