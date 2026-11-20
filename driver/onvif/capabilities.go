package onvif

import (
	"context"
	"fmt"

	usego "github.com/use-go/onvif"
	usedevice "github.com/use-go/onvif/device"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetCapabilities(ctx context.Context, camera Camera) (*model.Capabilities, error) {
	var result *model.Capabilities

	err := c.withDevice(ctx, camera, func(dev *usego.Device) error {
		resp, err := dev.CallMethod(usedevice.GetCapabilities{
			Category: "All",
		})
		if err != nil {
			return fmt.Errorf("get capabilities: %w", err)
		}

		result, err = xml.ParseCapabilities(resp)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
