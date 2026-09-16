package isapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetStatus(ctx context.Context, camera Camera) (*model.Status, error) {
	return c.getStatus(ctx, camera)
}

func (c *Client) getStatus(ctx context.Context, camera Camera) (*model.Status, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/status",
			camera.Channel,
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	return xml.ParseStatus(resp)
}
