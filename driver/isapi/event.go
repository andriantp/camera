package isapi

import (
	"context"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetEventCapabilities(ctx context.Context, camera Camera) (*model.EventCapability, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/Event/capabilities",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseEventCapabilities(resp)
}
