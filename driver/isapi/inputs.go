package isapi

import (
	"context"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetInputs(ctx context.Context, camera Camera) ([]model.Input, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/System/Video/inputs/channels",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseInputs(resp)
}
