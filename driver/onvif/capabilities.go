package onvif

import (
	"context"
	"fmt"

	usedevice "github.com/use-go/onvif/device"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetCapabilities(ctx context.Context, camera Camera) (*model.Capabilities, error) {
	dev, err := c.device(camera)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	resp, err := dev.CallMethod(usedevice.GetCapabilities{
		Category: "All",
	})
	if err != nil {
		return nil, fmt.Errorf("get capabilities: %w", err)
	}

	return xml.ParseCapabilities(resp)
}
