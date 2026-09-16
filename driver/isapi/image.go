package isapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetImageCapabilities(ctx context.Context, camera Camera) (*model.ImageCapabilities, error) {

	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf("/ISAPI/Image/channels/%d/capabilities", camera.Channel),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParseImageCapabilities(resp)
}
