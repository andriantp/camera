package isapi

import (
	"context"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetAnalyticsCapabilities(ctx context.Context, camera Camera) (*model.AnalyticsCapability, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/Smart/Capabilities",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseAnalyticsCapabilities(resp)
}
