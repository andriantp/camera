package onvif

import (
	"context"
	"fmt"

	usego "github.com/use-go/onvif"
	useptz "github.com/use-go/onvif/ptz"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetConfigurations(ctx context.Context, camera Camera) (*model.Configuration, error) {
	var result *model.Configuration

	err := c.withDevice(ctx, camera, func(dev *usego.Device) error {
		resp, err := dev.CallMethod(useptz.GetConfigurations{})
		if err != nil {
			return fmt.Errorf("get configurations: %w", err)
		}

		result, err = xml.ParseConfiguration(resp)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
