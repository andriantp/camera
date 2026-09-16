package isapi

import (
	"context"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetSystemCapabilities(ctx context.Context, camera Camera) (*model.SystemCapabilities, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/System/capabilities",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseSystemCapabilities(resp)
}
